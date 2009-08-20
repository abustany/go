// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go.h"

/*
 * machine size and rounding
 * alignment is dictated around
 * the size of a pointer, set in belexinit
 * (see ../6g/align.c).
 */

uint32
rnd(uint32 o, uint32 r)
{
	if(maxround == 0)
		fatal("rnd");

	if(r > maxround)
		r = maxround;
	if(r != 0)
		while(o%r != 0)
			o++;
	return o;
}

static void
offmod(Type *t)
{
	Type *f;
	int32 o;

	o = 0;
	for(f=t->type; f!=T; f=f->down) {
		if(f->etype != TFIELD)
			fatal("widstruct: not TFIELD: %lT", f);
		if(f->type->etype != TFUNC)
			continue;
		f->width = o;
		o += widthptr;
	}
}

static uint32
arrayelemwidth(Type *t)
{

	while(t->etype == TARRAY && t->bound >= 0)
		t = t->type;
	return t->width;
}

static uint32
widstruct(Type *t, uint32 o, int flag)
{
	Type *f;
	int32 w, m;

	for(f=t->type; f!=T; f=f->down) {
		if(f->etype != TFIELD)
			fatal("widstruct: not TFIELD: %lT", f);
		dowidth(f->type);
		if(f->type->width < 0 || f->type->width > 100000000)
			fatal("invalid width %lld", f->type->width);
		w = f->type->width;
		m = arrayelemwidth(f->type);
		o = rnd(o, m);
		f->width = o;	// really offset for TFIELD
		if(f->nname != N) {
			// this same stackparam logic is in addrescapes
			// in typecheck.c.  usually addrescapes runs after
			// widstruct, in which case we could drop this,
			// but function closure functions are the exception.
			if(f->nname->stackparam) {
				f->nname->stackparam->xoffset = o;
				f->nname->xoffset = 0;
			} else
				f->nname->xoffset = o;
		}
		o += w;
	}
	// final width is rounded
	if(flag)
		o = rnd(o, maxround);

	// type width only includes back to first field's offset
	if(t->type == T)
		t->width = 0;
	else
		t->width = o - t->type->width;
	return o;
}

void
dowidth(Type *t)
{
	int32 et;
	uint32 w;
	int lno;

	if(maxround == 0 || widthptr == 0)
		fatal("dowidth without betypeinit");

	if(t == T)
		return;

	if(t->width > 0)
		return;

	if(t->width == -2) {
		lno = lineno;
		lineno = t->lineno;
		yyerror("invalid recursive type %T", t);
		t->width = 0;
		lineno = lno;
		return;
	}

	lno = lineno;
	lineno = t->lineno;
	t->width = -2;

	et = t->etype;
	switch(et) {
	case TFUNC:
	case TCHAN:
	case TMAP:
	case TSTRING:
		break;

	default:
		/* simtype == 0 during bootstrap */
		if(simtype[t->etype] != 0)
			et = simtype[t->etype];
		break;
	}

	w = 0;
	switch(et) {
	default:
		fatal("dowidth: unknown type: %E", t->etype);
		break;

	/* compiler-specific stuff */
	case TINT8:
	case TUINT8:
	case TBOOL:		// bool is int8
		w = 1;
		break;
	case TINT16:
	case TUINT16:
		w = 2;
		break;
	case TINT32:
	case TUINT32:
	case TFLOAT32:
	case TPTR32:		// note lack of recursion
		w = 4;
		break;
	case TINT64:
	case TUINT64:
	case TFLOAT64:
	case TPTR64:		// note lack of recursion
		w = 8;
		break;
	case TFLOAT80:
		w = 10;
		break;
	case TDDD:
		w = 2*widthptr;
		break;
	case TINTER:		// implemented as 2 pointers
		offmod(t);
		w = 2*widthptr;
		break;
	case TCHAN:		// implemented as pointer
		dowidth(t->type);
		dowidth(t->down);
		w = widthptr;
		break;
	case TMAP:		// implemented as pointer
		dowidth(t->type);
		w = widthptr;
		break;
	case TFORW:		// should have been filled in
	case TANY:
		// dummy type; should be replaced before use.
		fatal("dowidth any");
		break;
	case TSTRING:
		if(sizeof_String == 0)
			fatal("early dowidth string");
		w = sizeof_String;
		break;
	case TARRAY:
		if(t->type == T)
			break;
		dowidth(t->type);
		w = sizeof_Array;
		if(t->bound >= 0)
			w = t->bound * t->type->width;
		break;

	case TSTRUCT:
		if(t->funarg)
			fatal("dowidth fn struct %T", t);
		w = widstruct(t, 0, 1);
		if(w == 0)
			w = maxround;
		break;

	case TFUNC:
		// function is 3 cated structures;
		// compute their widths as side-effect.
		w = widstruct(*getthis(t), 0, 0);
		w = widstruct(*getinarg(t), w, 1);
		w = widstruct(*getoutarg(t), w, 1);
		t->argwid = w;

		// but width of func type is pointer
		w = widthptr;
		break;
	}

	t->width = w;
	lineno = lno;
}

void
typeinit(void)
{
	int i, etype, sameas;
	Type *t;
	Sym *s;

	if(widthptr == 0)
		fatal("typeinit before betypeinit");

	for(i=0; i<NTYPE; i++)
		simtype[i] = i;

	types[TPTR32] = typ(TPTR32);
	dowidth(types[TPTR32]);

	types[TPTR64] = typ(TPTR64);
	dowidth(types[TPTR64]);

	tptr = TPTR32;
	if(widthptr == 8)
		tptr = TPTR64;

	for(i=TINT8; i<=TUINT64; i++)
		isint[i] = 1;
	isint[TINT] = 1;
	isint[TUINT] = 1;
	isint[TUINTPTR] = 1;

	for(i=TFLOAT32; i<=TFLOAT80; i++)
		isfloat[i] = 1;
	isfloat[TFLOAT] = 1;

	isptr[TPTR32] = 1;
	isptr[TPTR64] = 1;

	isforw[TFORW] = 1;

	issigned[TINT] = 1;
	issigned[TINT8] = 1;
	issigned[TINT16] = 1;
	issigned[TINT32] = 1;
	issigned[TINT64] = 1;

	/*
	 * initialize okfor
	 */
	for(i=0; i<NTYPE; i++) {
		if(isint[i] || i == TIDEAL) {
			okforeq[i] = 1;
			okforcmp[i] = 1;
			okforarith[i] = 1;
			okforadd[i] = 1;
			okforand[i] = 1;
			issimple[i] = 1;
			minintval[i] = mal(sizeof(*minintval[i]));
			maxintval[i] = mal(sizeof(*maxintval[i]));
		}
		if(isfloat[i]) {
			okforeq[i] = 1;
			okforcmp[i] = 1;
			okforadd[i] = 1;
			okforarith[i] = 1;
			issimple[i] = 1;
			minfltval[i] = mal(sizeof(*minfltval[i]));
			maxfltval[i] = mal(sizeof(*maxfltval[i]));
		}
	}

	issimple[TBOOL] = 1;

	okforadd[TSTRING] = 1;

	okforbool[TBOOL] = 1;

	okforcap[TARRAY] = 1;
	okforcap[TCHAN] = 1;

	okforlen[TARRAY] = 1;
	okforlen[TCHAN] = 1;
	okforlen[TMAP] = 1;
	okforlen[TSTRING] = 1;

	okforeq[TPTR32] = 1;
	okforeq[TPTR64] = 1;
	okforeq[TINTER] = 1;
	okforeq[TMAP] = 1;
	okforeq[TCHAN] = 1;
	okforeq[TFUNC] = 1;
	okforeq[TSTRING] = 1;
	okforeq[TBOOL] = 1;
	okforeq[TARRAY] = 1;	// refined in typecheck

	okforcmp[TSTRING] = 1;

	for(i=0; i<nelem(okfor); i++)
		okfor[i] = okfornone;

	// binary
	okfor[OADD] = okforadd;
	okfor[OAND] = okforand;
	okfor[OANDAND] = okforbool;
	okfor[OANDNOT] = okforand;
	okfor[ODIV] = okforarith;
	okfor[OEQ] = okforeq;
	okfor[OGE] = okforcmp;
	okfor[OGT] = okforcmp;
	okfor[OLE] = okforcmp;
	okfor[OLT] = okforcmp;
	okfor[OMOD] = okforarith;
	okfor[OMUL] = okforarith;
	okfor[ONE] = okforeq;
	okfor[OOR] = okforand;
	okfor[OOROR] = okforbool;
	okfor[OSUB] = okforarith;
	okfor[OXOR] = okforand;
	okfor[OLSH] = okforand;
	okfor[ORSH] = okforand;

	// unary
	okfor[OCOM] = okforand;
	okfor[OMINUS] = okforarith;
	okfor[ONOT] = okforbool;
	okfor[OPLUS] = okforadd;

	// special
	okfor[OCAP] = okforcap;
	okfor[OLEN] = okforlen;

	// comparison
	iscmp[OLT] = 1;
	iscmp[OGT] = 1;
	iscmp[OGE] = 1;
	iscmp[OLE] = 1;
	iscmp[OEQ] = 1;
	iscmp[ONE] = 1;

	mpatofix(maxintval[TINT8], "0x7f");
	mpatofix(minintval[TINT8], "-0x80");
	mpatofix(maxintval[TINT16], "0x7fff");
	mpatofix(minintval[TINT16], "-0x8000");
	mpatofix(maxintval[TINT32], "0x7fffffff");
	mpatofix(minintval[TINT32], "-0x80000000");
	mpatofix(maxintval[TINT64], "0x7fffffffffffffff");
	mpatofix(minintval[TINT64], "-0x8000000000000000");

	mpatofix(maxintval[TUINT8], "0xff");
	mpatofix(maxintval[TUINT16], "0xffff");
	mpatofix(maxintval[TUINT32], "0xffffffff");
	mpatofix(maxintval[TUINT64], "0xffffffffffffffff");

	mpatoflt(maxfltval[TFLOAT32], "3.40282347e+38");
	mpatoflt(minfltval[TFLOAT32], "-3.40282347e+38");
	mpatoflt(maxfltval[TFLOAT64], "1.7976931348623157e+308");
	mpatoflt(minfltval[TFLOAT64], "-1.7976931348623157e+308");

	/* for walk to use in error messages */
	types[TFUNC] = functype(N, nil, nil);

	/* types used in front end */
	// types[TNIL] got set early in lexinit
	types[TIDEAL] = typ(TIDEAL);

	/* simple aliases */
	simtype[TMAP] = tptr;
	simtype[TCHAN] = tptr;
	simtype[TFUNC] = tptr;

	/* pick up the backend typedefs */
	for(i=0; typedefs[i].name; i++) {
		s = lookup(typedefs[i].name);

		etype = typedefs[i].etype;
		if(etype < 0 || etype >= nelem(types))
			fatal("typeinit: %s bad etype", s->name);
		sameas = typedefs[i].sameas;
		if(sameas < 0 || sameas >= nelem(types))
			fatal("typeinit: %s bad sameas", s->name);
		simtype[etype] = sameas;
		minfltval[etype] = minfltval[sameas];
		maxfltval[etype] = maxfltval[sameas];
		minintval[etype] = minintval[sameas];
		maxintval[etype] = maxintval[sameas];

		t = types[etype];
		if(t != T)
			fatal("typeinit: %s already defined", s->name);

		t = typ(etype);
		t->sym = s;

		dowidth(t);
		types[etype] = t;
		s->def = typenod(t);
	}

	Array_array = rnd(0, widthptr);
	Array_nel = rnd(Array_array+widthptr, types[TUINT32]->width);
	Array_cap = rnd(Array_nel+types[TUINT32]->width, types[TUINT32]->width);
	sizeof_Array = rnd(Array_cap+types[TUINT32]->width, maxround);

	// string is same as slice wo the cap
	sizeof_String = rnd(Array_nel+types[TUINT32]->width, maxround);

	dowidth(types[TSTRING]);
	dowidth(idealstring);
}

/*
 * compute total size of f's in/out arguments.
 */
int
argsize(Type *t)
{
	Iter save;
	Type *fp;
	int w, x;

	w = 0;

	fp = structfirst(&save, getoutarg(t));
	while(fp != T) {
		x = fp->width + fp->type->width;
		if(x > w)
			w = x;
		fp = structnext(&save);
	}

	fp = funcfirst(&save, t);
	while(fp != T) {
		x = fp->width + fp->type->width;
		if(x > w)
			w = x;
		fp = funcnext(&save);
	}

	w = (w+7) & ~7;
	return w;
}
