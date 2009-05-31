// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%{
#include "go.h"
%}
%union	{
	Node*		node;
	Type*		type;
	Sym*		sym;
	struct	Val	val;
	int		lint;
}
%token	<val>		LLITERAL
%token	<lint>		LASOP
%token	<sym>		LNAME LATYPE LPACK
%token	<sym>		LPACKAGE LIMPORT LDEFER LCLOSE LCLOSED
%token	<sym>		LMAP LCHAN LINTERFACE LFUNC LSTRUCT
%token	<sym>		LCOLAS LFALL LRETURN LDDD
%token	<sym>		LLEN LCAP LPANIC LPANICN LPRINT LPRINTN
%token	<sym>		LVAR LTYPE LCONST LSELECT LMAKE LNEW
%token	<sym>		LFOR LIF LELSE LSWITCH LCASE LDEFAULT
%token	<sym>		LBREAK LCONTINUE LGO LGOTO LRANGE
%token	<sym>		LNIL LTRUE LFALSE LIOTA

%token			LOROR LANDAND LEQ LNE LLE LLT LGE LGT
%token			LLSH LRSH LINC LDEC LCOMM LANDNOT
%token			LIGNORE

/*
 * the go semicolon rules are:
 *
 *  1. all statements and declarations are terminated by semicolons
 *  2. semicolons can be omitted at top level.
 *  3. semicolons can be omitted before and after the closing ) or }
 *	on a list of statements or declarations.
 *
 * thus the grammar must distinguish productions that
 * can omit the semicolon terminator and those that can't.
 * names like Astmt, Avardcl, etc. can drop the semicolon.
 * names like Bstmt, Bvardcl, etc. can't.
 */

%type	<sym>		sym sym1 sym2 sym3 keyword lname latype lpackatype
%type	<node>		xdcl xdcl_list_r oxdcl_list
%type	<node>		common_dcl Acommon_dcl Bcommon_dcl
%type	<node>		oarg_type_list arg_type_list_r arg_chunk arg_chunk_list_r arg_type_list
%type	<node>		Aelse_stmt Belse_stmt
%type	<node>		complex_stmt compound_stmt switch_body ocaseblock_list ostmt_list
%type	<node>		caseblock_list_r stmt_list_r Astmt_list_r Bstmt_list_r
%type	<node>		Astmt Bstmt
%type	<node>		for_stmt for_body for_header
%type	<node>		if_stmt if_header select_stmt switch_stmt condition case caseblock
%type	<node>		simple_stmt osimple_stmt range_stmt semi_stmt
%type	<node>		expr uexpr pexpr expr_list oexpr oexpr_list expr_list_r
%type	<node>		exprsym3_list_r exprsym3 pseudocall
%type	<node>		name labelname onew_name new_name new_name_list_r new_field
%type	<node>		vardcl_list_r vardcl Avardcl Bvardcl
%type	<node>		interfacedcl_list_r interfacedcl interfacedcl1
%type	<node>		structdcl_list_r structdcl embed
%type	<node>		fnres Afnres Bfnres fnliteral xfndcl fndcl fnbody
%type	<node>		braced_keyexpr_list keyval_list_r keyval

%type	<type>		typedclname new_type
%type	<type>		type Atype Btype
%type	<type>		othertype Aothertype Bothertype
%type	<type>		chantype Achantype Bchantype
%type	<type>		fntype Afntype Bfntype
%type	<type>		nametype structtype interfacetype convtype
%type	<type>		non_name_type Anon_fn_type Bnon_fn_type
%type	<type>		Anon_chan_type Bnon_chan_type
%type	<type>		indcl fnlitdcl dotdotdot
%type	<val>		oliteral

%type	<node>		hidden_constant
%type	<node>		hidden_dcl hidden_structdcl
%type	<type>		hidden_type hidden_type1 hidden_type2
%type	<node>		hidden_structdcl_list ohidden_structdcl_list hidden_structdcl_list_r
%type	<node>		hidden_interfacedcl_list ohidden_interfacedcl_list hidden_interfacedcl_list_r
%type	<node>		hidden_interfacedcl
%type	<node>		hidden_funarg_list ohidden_funarg_list hidden_funarg_list_r
%type	<node>		hidden_funres ohidden_funres hidden_importsym hidden_pkg_importsym

%left			LOROR
%left			LANDAND
%left			LCOMM
%left			LEQ LNE LLE LGE LLT LGT
%left			'+' '-' '|' '^'
%left			'*' '/' '%' '&' LLSH LRSH LANDNOT

/*
 * resolve { vs condition in favor of condition
 */
%left			'{'
%left			Condition

/*
 * resolve LPACKAGE vs not in favor of LPACKAGE
 */
%left			NotPackage
%left			LPACKAGE

/*
 * resolve '.' vs not in favor of '.'
 */
%left			NotDot
%left			'.'

/*
 * resolve '(' vs not in favor of '('
 */
%left			NotParen
%left			'('

%%
file:
	loadsys
	package
	imports
	oxdcl_list
	{
		if(debug['f'])
			frame(1);
		fninit($4);
		testdclstack();
	}

package:
	%prec NotPackage
	{
		yyerror("package statement must be first");
		mkpackage("main");
	}
|	LPACKAGE sym
	{
		mkpackage($2->name);
	}

/*
 * this loads the definitions for the sys functions,
 * so that the compiler can generate calls to them,
 * but does not make the name "sys" visible as a package.
 */
loadsys:
	{
		cannedimports("sys.6", sysimport);
	}
	import_package
	import_there
	{
		pkgimportname = S;
	}

imports:
|	imports import

import:
	LIMPORT import_stmt
|	LIMPORT '(' import_stmt_list_r osemi ')'
|	LIMPORT '(' ')'

import_stmt:
	import_here import_package import_there import_done

import_here:
	LLITERAL
	{
		// import with original name
		pkgimportname = S;
		pkgmyname = S;
		importfile(&$1);
	}
|	sym LLITERAL
	{
		// import with given name
		pkgimportname = S;
		pkgmyname = $1;
		importfile(&$2);
	}
|	'.' LLITERAL
	{
		// import into my name space
		pkgmyname = lookup(".");
		importfile(&$2);
	}

import_package:
	LPACKAGE sym
	{
		pkgimportname = $2;
		if(strcmp($2->name, "main") == 0)
			yyerror("cannot import package main");
	}

import_there:
	hidden_import_list '$' '$'
	{
		checkimports();
		unimportfile();
	}
|	LIMPORT '$' '$' hidden_import_list '$' '$'
	{
		checkimports();
	}

import_done:
	{
		Sym *import, *my;

		import = pkgimportname;
		my = pkgmyname;
		pkgmyname = S;
		pkgimportname = S;

		if(import == S)
			break;
		if(my == S)
			my = import;
		if(my->name[0] == '.') {
			importdot(import);
			break;
		}

		// In order to allow multifile packages to use type names
		// that are the same as the package name (i.e. go/parser
		// is package parser and has a type called parser), we have
		// to not bother trying to declare the package if it is our package.
		// TODO(rsc): Is there a better way to tell if the package is ours?
		if(my == import && strcmp(import->name, package) == 0)
			break;

		if(my->lexical != LNAME || my->oname != N || my->otype != T) {
			// TODO(rsc): this line is only needed because of the
			//	package net
			//	import "net"
			// convention; if we get rid of it, the check can go away
			// and we can just always print the error
			if(my->lexical != LPACK || strcmp(my->opack, import->name) != 0)
				yyerror("redeclaration of %S by import", my);
		}
		my->lexical = LPACK;
		my->opack = import->name;
	}

hidden_import_list:
	{
		defercheckwidth();
	}
	hidden_import_list_r
	{
		resumecheckwidth();
	}

/*
 * declarations
 */
xdcl:
	{ stksize = initstksize; } common_dcl
	{
		$$ = $2;
		initstksize = stksize;
	}
|	xfndcl
	{
		if($1 != N && $1->nname != N && $1->type->thistuple == 0)
			autoexport($1->nname->sym);
		$$ = N;
	}
|	';'
	{
		$$ = N;
	}
|	error xdcl
	{
		$$ = $2;
	}

common_dcl:
	Acommon_dcl
|	Bcommon_dcl

Acommon_dcl:
	LVAR Avardcl
	{
		$$ = $2;
	}
|	LVAR '(' vardcl_list_r osemi ')'
	{
		$$ = rev($3);
	}
|	LVAR '(' ')'
	{
		$$ = N;
	}
|	LCONST '(' constdcl osemi ')'
	{
		iota = 0;
		lastconst = N;
		$$ = N;
	}
|	LCONST '(' constdcl ';' constdcl_list_r osemi ')'
	{
		iota = 0;
		lastconst = N;
		$$ = N;
	}
|	LCONST '(' ')'
	{
		$$ = N;
	}
|	LTYPE Atypedcl
	{
		$$ = N;
	}
|	LTYPE '(' typedcl_list_r osemi ')'
	{
		$$ = N;
	}
|	LTYPE '(' ')'
	{
		$$ = N;
	}

Bcommon_dcl:
	LVAR Bvardcl
	{
		$$ = $2;
	}
|	LCONST constdcl
	{
		$$ = N;
		iota = 0;
		lastconst = N;
	}
|	LTYPE Btypedcl
	{
		$$ = N;
	}

vardcl:
	Avardcl
|	Bvardcl

Avardcl:
	new_name_list_r Atype
	{
		$$ = rev($1);
		dodclvar($$, $2);

		if(funcdepth == 0) {
			$$ = N;
		} else {
			$$ = nod(OAS, $$, N);
			addtotop($$);
		}
	}

Bvardcl:
	new_name_list_r Btype
	{
		$$ = rev($1);
		dodclvar($$, $2);

		if(funcdepth == 0) {
			$$ = N;
		} else {
			$$ = nod(OAS, $$, N);
			addtotop($$);
		}
	}
|	new_name_list_r type '=' expr_list
	{
		if(addtop != N)
			fatal("new_name_list_r type '=' expr_list");

		$$ = variter($1, $2, $4);
		addtotop($$);
	}
|	new_name_list_r '=' expr_list
	{
		if(addtop != N)
			fatal("new_name_list_r '=' expr_list");

		$$ = variter($1, T, $3);
		addtotop($$);
	}

constdcl:
	new_name_list_r type '=' expr_list
	{
		constiter($1, $2, $4);
	}
|	new_name_list_r '=' expr_list
	{
		constiter($1, T, $3);
	}

constdcl1:
	constdcl
|	new_name_list_r type
	{
		constiter($1, $2, N);
	}
|	new_name_list_r
	{
		constiter($1, T, N);
	}

typedclname:
	new_type
	{
		$$ = dodcltype($1);
		defercheckwidth();
	}

typedcl:
	Atypedcl
|	Btypedcl

Atypedcl:
	typedclname Atype
	{
		updatetype($1, $2);
		resumecheckwidth();
	}

Btypedcl:
	typedclname Btype
	{
		updatetype($1, $2);
		resumecheckwidth();
	}
|	typedclname LSTRUCT
	{
		updatetype($1, typ(TFORWSTRUCT));
		resumecheckwidth();
	}
|	typedclname LINTERFACE
	{
		updatetype($1, typ(TFORWINTER));
		resumecheckwidth();
	}

Aelse_stmt:
	complex_stmt
|	compound_stmt

Belse_stmt:
	simple_stmt
|	semi_stmt
|	';'
	{
		$$ = N;
	}

simple_stmt:
	expr
	{
		$$ = $1;
	}
|	expr LASOP expr
	{
		$$ = nod(OASOP, $1, $3);
		$$->etype = $2;			// rathole to pass opcode
	}
|	exprsym3_list_r '=' expr_list
	{
		$$ = rev($1);
		$$ = nod(OAS, $$, $3);
	}
|	exprsym3_list_r LCOLAS expr_list
	{
		if(addtop != N)
			fatal("exprsym3_list_r LCOLAS expr_list");
		if($3->op == OTYPESW) {
			$$ = nod(OTYPESW, $1, $3->left);
			break;
		}
		$$ = rev($1);
		$$ = colas($$, $3);
		$$ = nod(OAS, $$, $3);
		$$->colas = 1;
		addtotop($$);
	}
|	LPRINT '(' oexpr_list ')'
	{
		$$ = nod(OPRINT, $3, N);
	}
|	LPRINTN '(' oexpr_list ')'
	{
		$$ = nod(OPRINTN, $3, N);
	}
|	LPANIC '(' oexpr_list ')'
	{
		$$ = nod(OPANIC, $3, N);
	}
|	LPANICN '(' oexpr_list ')'
	{
		$$ = nod(OPANICN, $3, N);
	}
|	expr LINC
	{
		$$ = nod(OASOP, $1, nodintconst(1));
		$$->etype = OADD;
	}
|	expr LDEC
	{
		$$ = nod(OASOP, $1, nodintconst(1));
		$$->etype = OSUB;
	}

complex_stmt:
	for_stmt
|	switch_stmt
|	select_stmt
|	if_stmt
	{
		popdcl();
		$$ = $1;
	}
|	if_stmt LELSE Aelse_stmt
	{
		popdcl();
		$$ = $1;
		$$->nelse = $3;
	}

case:
	LCASE expr_list ':'
	{
		// will be converted to OCASE
		// right will point to next case
		// done in casebody()
		poptodcl();
		if(typeswvar != N && typeswvar->right != N)
		if($2->op == OLITERAL && $2->val.ctype == CTNIL) {
			// this version in type switch case nil
			$$ = nod(OTYPESW, N, N);
			$$ = nod(OXCASE, $$, N);
			break;
		}
		$$ = nod(OXCASE, $2, N);
	}
|	LCASE name '=' expr ':'
	{
		// will be converted to OCASE
		// right will point to next case
		// done in casebody()
		poptodcl();
		$$ = nod(OAS, $2, $4);
		$$ = nod(OXCASE, $$, N);
	}
|	LCASE name LCOLAS expr ':'
	{
		// will be converted to OCASE
		// right will point to next case
		// done in casebody()
		poptodcl();
		$$ = nod(OAS, selectas($2,$4), $4);
		$$ = nod(OXCASE, $$, N);
		addtotop($$);
	}
|	LCASE type ':'
	{
		poptodcl();
		if(typeswvar == N || typeswvar->right == N) {
			yyerror("type case not in a type switch");
			$$ = N;
		} else
			$$ = old2new(typeswvar->right, $2);
		$$ = nod(OTYPESW, $$, N);
		$$ = nod(OXCASE, $$, N);
		addtotop($$);
	}
|	LDEFAULT ':'
	{
		poptodcl();
		$$ = nod(OXCASE, N, N);
	}

semi_stmt:
	LFALL
	{
		// will be converted to OFALL
		$$ = nod(OXFALL, N, N);
	}
|	LBREAK onew_name
	{
		$$ = nod(OBREAK, $2, N);
	}
|	LCONTINUE onew_name
	{
		$$ = nod(OCONTINUE, $2, N);
	}
|	LGO pseudocall
	{
		$$ = nod(OPROC, $2, N);
	}
|	LDEFER pseudocall
	{
		$$ = nod(ODEFER, $2, N);
	}
|	LGOTO new_name
	{
		$$ = nod(OGOTO, $2, N);
	}
|	LRETURN oexpr_list
	{
		$$ = nod(ORETURN, $2, N);
	}
|	if_stmt LELSE Belse_stmt
	{
		popdcl();
		$$ = $1;
		$$->nelse = $3;
	}

compound_stmt:
	'{'
	{
		markdcl();
	} ostmt_list '}'
	{
		$$ = $3;
		if($$ == N)
			$$ = nod(OEMPTY, N, N);
		popdcl();
	}

switch_body:
	'{'
	{
		markdcl();
	}
	ocaseblock_list '}'
	{
		$$ = $3;
		if($$ == N)
			$$ = nod(OEMPTY, N, N);
		popdcl();
	}

caseblock:
	case ostmt_list
	{
		$$ = $1;
		$$->nbody = $2;
	}

caseblock_list_r:
	caseblock
|	caseblock_list_r caseblock
	{
		$$ = nod(OLIST, $1, $2);
	}

range_stmt:
	exprsym3_list_r '=' LRANGE expr
	{
		$$ = nod(ORANGE, $1, $4);
		$$->etype = 0;	// := flag
	}
|	exprsym3_list_r LCOLAS LRANGE expr
	{
		$$ = nod(ORANGE, $1, $4);
		$$->etype = 1;
	}

for_header:
	osimple_stmt ';' osimple_stmt ';' osimple_stmt
	{
		// init ; test ; incr
		if($5 != N && $5->colas != 0)
			yyerror("cannot declare in the for-increment");
		$$ = nod(OFOR, N, N);
		$$->ninit = $1;
		$$->ntest = $3;
		$$->nincr = $5;
	}
|	condition
	{
		// normal test
		$$ = nod(OFOR, N, N);
		$$->ninit = N;
		$$->ntest = $1;
		$$->nincr = N;
	}
|	range_stmt
	{
		$$ = dorange($1);
		addtotop($$);
	}

for_body:
	for_header compound_stmt
	{
		$$ = $1;
		$$->nbody = list($$->nbody, $2);
	}

for_stmt:
	LFOR
	{
		markdcl();
	}
	for_body
	{
		$$ = $3;
		popdcl();
	}

/*
 * using cond instead of osimple_stmt creates
 * a shift/reduce conflict on an input like
 *
 *	if x == []int { true } { true }
 *
 * at the first {, giving us an opportunity
 * to resolve it by reduce, which implements
 * the rule about { } inside if conditions
 * needing parens.
 */
condition:
	osimple_stmt	%prec Condition


if_header:
	condition
	{
		// test
		$$ = nod(OIF, N, N);
		$$->ninit = N;
		$$->ntest = $1;
	}
|	osimple_stmt ';' condition
	{
		// init ; test
		$$ = nod(OIF, N, N);
		$$->ninit = $1;
		$$->ntest = $3;
	}

if_stmt:
	LIF
	{
		markdcl();
	}
	if_header compound_stmt
	{
		$$ = $3;
		$$->nbody = $4;
		// no popdcl; maybe there's an LELSE
	}

switch_stmt:
	LSWITCH
	{
		markdcl();
	}
	if_header
	{
		Node *n;
		n = $3->ntest;
		if(n != N && n->op == OTYPESW)
			n = n->left;
		else
			n = N;
		typeswvar = nod(OLIST, typeswvar, n);
	}
	switch_body
	{
		$$ = $3;
		$$->op = OSWITCH;
		$$->nbody = $5;
		typeswvar = typeswvar->left;
		popdcl();
	}

select_stmt:
	LSELECT
	{
		markdcl();
	}
	switch_body
	{
		$$ = nod(OSELECT, $3, N);
		popdcl();
	}

/*
 * expressions
 */
expr:
	uexpr
|	expr LOROR expr
	{
		$$ = nod(OOROR, $1, $3);
	}
|	expr LANDAND expr
	{
		$$ = nod(OANDAND, $1, $3);
	}
|	expr LEQ expr
	{
		$$ = nod(OEQ, $1, $3);
	}
|	expr LNE expr
	{
		$$ = nod(ONE, $1, $3);
	}
|	expr LLT expr
	{
		$$ = nod(OLT, $1, $3);
	}
|	expr LLE expr
	{
		$$ = nod(OLE, $1, $3);
	}
|	expr LGE expr
	{
		$$ = nod(OGE, $1, $3);
	}
|	expr LGT expr
	{
		$$ = nod(OGT, $1, $3);
	}
|	expr '+' expr
	{
		$$ = nod(OADD, $1, $3);
	}
|	expr '-' expr
	{
		$$ = nod(OSUB, $1, $3);
	}
|	expr '|' expr
	{
		$$ = nod(OOR, $1, $3);
	}
|	expr '^' expr
	{
		$$ = nod(OXOR, $1, $3);
	}
|	expr '*' expr
	{
		$$ = nod(OMUL, $1, $3);
	}
|	expr '/' expr
	{
		$$ = nod(ODIV, $1, $3);
	}
|	expr '%' expr
	{
		$$ = nod(OMOD, $1, $3);
	}
|	expr '&' expr
	{
		$$ = nod(OAND, $1, $3);
	}
|	expr LANDNOT expr
	{
		$$ = nod(OANDNOT, $1, $3);
	}
|	expr LLSH expr
	{
		$$ = nod(OLSH, $1, $3);
	}
|	expr LRSH expr
	{
		$$ = nod(ORSH, $1, $3);
	}
|	expr LCOMM expr
	{
		$$ = nod(OSEND, $1, $3);
	}

uexpr:
	pexpr
|	'*' uexpr
	{
		$$ = nod(OIND, $2, N);
	}
|	'&' uexpr
	{
		$$ = nod(OADDR, $2, N);
	}
|	'+' uexpr
	{
		$$ = nod(OPLUS, $2, N);
	}
|	'-' uexpr
	{
		$$ = nod(OMINUS, $2, N);
	}
|	'!' uexpr
	{
		$$ = nod(ONOT, $2, N);
	}
|	'~' uexpr
	{
		yyerror("the OCOM operator is ^");
		$$ = nod(OCOM, $2, N);
	}
|	'^' uexpr
	{
		$$ = nod(OCOM, $2, N);
	}
|	LCOMM uexpr
	{
		$$ = nod(ORECV, $2, N);
	}

/*
 * call-like statements that
 * can be preceeded by 'defer' and 'go'
 */
pseudocall:
	pexpr '(' oexpr_list ')'
	{
		$$ = unsafenmagic($1, $3);
		if($$ == N)
			$$ = nod(OCALL, $1, $3);
	}
|	LCLOSE '(' expr ')'
	{
		$$ = nod(OCLOSE, $3, N);
	}

pexpr:
	LLITERAL
	{
		$$ = nodlit($1);
	}
|	LNIL
	{
		Val v;
		v.ctype = CTNIL;
		$$ = nodlit(v);
	}
|	LTRUE
	{
		$$ = nodbool(1);
	}
|	LFALSE
	{
		$$ = nodbool(0);
	}
|	LIOTA
	{
		$$ = nodintconst(iota);
		$$->iota = 1;	// flag to reevaluate on copy
	}
|	name
|	'(' expr ')'
	{
		$$ = $2;
	}
|	pexpr '.' sym2
	{
		$$ = nod(ODOT, $1, newname($3));
		$$ = adddot($$);
	}
|	pexpr '.' '(' type ')'
	{
		$$ = nod(ODOTTYPE, $1, N);
		$$->type = $4;
	}
|	pexpr '.' '(' LTYPE ')'
	{
		$$ = nod(OTYPESW, $1, N);
	}
|	pexpr '[' expr ']'
	{
		$$ = nod(OINDEX, $1, $3);
	}
|	pexpr '[' keyval ']'
	{
		$$ = nod(OSLICE, $1, $3);
	}
|	pseudocall
|	LLEN '(' expr ')'
	{
		$$ = nod(OLEN, $3, N);
	}
|	LCLOSED '(' expr ')'
	{
		$$ = nod(OCLOSED, $3, N);
	}
|	LCAP '(' expr ')'
	{
		$$ = nod(OCAP, $3, N);
	}
|	LNEW '(' type ')'
	{
		$$ = nod(ONEW, N, N);
		$$->type = $3;
	}
|	LNEW '(' type ',' expr_list ')'
	{
		$$ = nod(ONEW, $5, N);
		$$->type = $3;
	}
|	LMAKE '(' type ')'
	{
		$$ = nod(OMAKE, N, N);
		$$->type = $3;
	}
|	LMAKE '(' type ',' expr_list ')'
	{
		$$ = nod(OMAKE, $5, N);
		$$->type = $3;
	}
|	convtype '(' expr ')'
	{
		// conversion
		$$ = rev($3);
		if($$ == N)
			$$ = nod(OEMPTY, N, N);
		$$ = nod(OCONV, $$, N);
		$$->type = $1;
	}
|	convtype '{' braced_keyexpr_list '}'
	{
		// composite expression
		$$ = rev($3);
		if($$ == N)
			$$ = nod(OEMPTY, N, N);
		$$ = nod(OCOMPOS, $$, N);
		$$->type = $1;
	}
|	fnliteral

/*
 * lexical symbols that can be
 * from other packages
 */
lpack:
	LPACK
	{
		context = $1->opack;
	}
/*
 * adding this would enable gri's nested package idea
 *
|	lpack '.' LPACK
	{
		context = $3->opack;
	}
 */

lname:
	LNAME
|	lpack '.' LNAME
	{
		$$ = $3;
		context = nil;
	}

latype:
	LATYPE
|	lpackatype

lpackatype:
	lpack '.' LATYPE
	{
		$$ = $3;
		context = nil;
	}

/*
 * names and types
 *	newname is used before declared
 *	oldname is used after declared
 */
new_name:
	sym1
	{
		$$ = newname($1);
	}

new_field:
	sym2
	{
		$$ = newname($1);
	}

new_type:
	sym1
	{
		$$ = newtype($1);
	}

onew_name:
	{
		$$ = N;
	}
|	new_name

sym:
	LATYPE
|	LNAME
|	LPACK

sym1:
	sym
|	keyword

/*
 * keywords that can be field names
 * pretty much any name can be allowed
 * limited only by good taste
 */
sym2:
	sym1

/*
 * keywords that can be variables
 * but are not already legal expressions
 */
sym3:
	LLEN
|	LCAP
|	LCLOSE
|	LCLOSED
|	LPANIC
|	LPANICN
|	LPRINT
|	LPRINTN
|	LNEW
|	LMAKE

/*
 * keywords that we can
 * use as variable/type names
 */
keyword:
	sym3
|	LNIL
|	LTRUE
|	LFALSE
|	LIOTA

name:
	lname
	{
		$$ = oldname($1);
	}
	/*
	 * this rule introduces 1 reduce/reduce conflict
	 * with the rule lpack: LPACK above.
	 * the reduce/reduce conflict is only with
	 * lookahead '.', in which case the correct
	 * resolution is the lpack rule.  (and it wins
	 * because it is above.)
	 */
|	LPACK	%prec NotDot
	{
		$$ = oldname($1);
	}

labelname:
	name
|	LATYPE
	{
		$$ = oldname($1);
	}
|	keyword
	{
		$$ = oldname($1);
	}

convtype:
	latype
	{
		$$ = oldtype($1);
	}
|	'[' oexpr ']' type
	{
		// array literal
		$$ = aindex($2, $4);
	}
|	'[' LDDD ']' type
	{
		// array literal of nelem
		$$ = aindex(N, $4);
		$$->bound = -100;
	}
|	LMAP '[' type ']' type
	{
		// map literal
		$$ = maptype($3, $5);
	}
|	structtype
|	'(' type ')'
	{
		$$ = $2;
	}

/*
 * to avoid parsing conflicts, type is split into
 *	named types
 *	channel types
 *	function types
 *	any other type
 *
 * (and also into A/B as described above).
 *
 * the type system makes additional restrictions,
 * but those are not implemented in the grammar.
 */
type:
	Atype
|	Btype

Atype:
	Achantype
|	Afntype
|	Aothertype

Btype:
	nametype
|	Bchantype
|	Bfntype
|	Bothertype
|	'(' type ')'
	{
		$$ = $2;
	}

non_name_type:
	chantype
|	fntype
|	othertype
|	dotdotdot

dotdotdot:
	LDDD
	{
		$$ = typ(TDDD);
	}

Anon_chan_type:
	Afntype
|	Aothertype

Bnon_chan_type:
	nametype
|	Bfntype
|	Bothertype
|	'(' Btype ')'
	{
		$$ = $2;
	}

Anon_fn_type:
	Achantype
|	Aothertype

Bnon_fn_type:
	nametype
|	Bchantype
|	Bothertype

nametype:
	LATYPE
	{
		if($1->otype != T && $1->otype->etype == TANY)
		if(strcmp(package, "PACKAGE") != 0)
			yyerror("the any type is restricted");
		$$ = oldtype($1);
	}

othertype:
	Aothertype
|	Bothertype

Aothertype:
	'[' oexpr ']' Atype
	{
		$$ = aindex($2, $4);
	}
|	LCOMM LCHAN Atype
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Crecv;
	}
|	LCHAN LCOMM Anon_chan_type
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Csend;
	}
|	LMAP '[' type ']' Atype
	{
		$$ = maptype($3, $5);
	}
|	'*' Atype
	{
		$$ = ptrto($2);
	}
|	structtype
|	interfacetype

Bothertype:
	lpackatype
	{
		$$ = oldtype($1);
	}
|	'[' oexpr ']' Btype
	{
		$$ = aindex($2, $4);
	}
|	LCOMM LCHAN Btype
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Crecv;
	}
|	LCHAN LCOMM Bnon_chan_type
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Csend;
	}
|	LMAP '[' type ']' Btype
	{
		$$ = maptype($3, $5);
	}
|	'*' Btype
	{
		$$ = ptrto($2);
	}

chantype:
	Achantype
|	Bchantype

Achantype:
	LCHAN Atype
	{
		$$ = typ(TCHAN);
		$$->type = $2;
		$$->chan = Cboth;
	}

Bchantype:
	LCHAN Btype
	{
		$$ = typ(TCHAN);
		$$->type = $2;
		$$->chan = Cboth;
	}

structtype:
	LSTRUCT '{' structdcl_list_r osemi '}'
	{
		$$ = dostruct(rev($3), TSTRUCT);
	}
|	LSTRUCT '{' '}'
	{
		$$ = dostruct(N, TSTRUCT);
	}

interfacetype:
	LINTERFACE '{' interfacedcl_list_r osemi '}'
	{
		$$ = dostruct(rev($3), TINTER);
		$$ = sortinter($$);
	}
|	LINTERFACE '{' '}'
	{
		$$ = dostruct(N, TINTER);
	}

keyval:
	expr ':' expr
	{
		$$ = nod(OKEY, $1, $3);
	}
|	LATYPE ':' expr
	{
		$$ = nod(OKEY, newname($1), $3);
	}


/*
 * function stuff
 * all in one place to show how crappy it all is
 */
xfndcl:
	LFUNC
	{
		maxarg = 0;
		stksize = 0;
	} fndcl fnbody
	{
		$$ = $3;
		$$->nbody = $4;
		funcbody($$);
	}

fndcl:
	new_name '(' oarg_type_list ')' fnres
	{
		b0stack = dclstack;	// mark base for fn literals
		$$ = nod(ODCLFUNC, N, N);
		$$->nname = $1;
		if($3 == N && $5 == N)
			$$->nname = renameinit($1);
		$$->type = functype(N, $3, $5);
		funchdr($$);
	}
|	'(' oarg_type_list ')' new_name '(' oarg_type_list ')' fnres
	{
		b0stack = dclstack;	// mark base for fn literals
		$$ = nod(ODCLFUNC, N, N);
		if(listcount($2) == 1) {
			$$->nname = $4;
			$$->nname = methodname($4, $2->type);
			$$->type = functype($2, $6, $8);
			funchdr($$);
			addmethod($4, $$->type, 1);
		} else {
			/* declare it as a function */
			yyerror("unknown method receiver");
			$$->nname = $4;
			$$->type = functype(N, $6, $8);
			funchdr($$);
		}

	}

fntype:
	Afntype
|	Bfntype

Afntype:
	LFUNC '(' oarg_type_list ')' Afnres
	{
		$$ = functype(N, $3, $5);
	}

Bfntype:
	LFUNC '(' oarg_type_list ')' Bfnres
	{
		$$ = functype(N, $3, $5);
	}

fnlitdcl:
	fntype
	{
		markdcl();
		$$ = $1;
		funclit0($$);
	}

fnliteral:
	fnlitdcl '{' ostmt_list '}'
	{
		$$ = funclit1($1, $3);
	}

fnbody:
	'{' ostmt_list '}'
	{
		$$ = $2;
		if($$ == N)
			$$ = nod(ORETURN, N, N);
	}
|	{
		$$ = N;
	}

fnres:
	Afnres
|	Bfnres

Afnres:
	Anon_fn_type
	{
		$$ = nod(ODCLFIELD, N, N);
		$$->type = $1;
		$$ = cleanidlist($$);
	}

Bfnres:
	%prec NotParen
	{
		$$ = N;
	}
|	Bnon_fn_type
	{
		$$ = nod(ODCLFIELD, N, N);
		$$->type = $1;
		$$ = cleanidlist($$);
	}
|	'(' oarg_type_list ')'
	{
		$$ = $2;
	}

/*
 * lists of things
 * note that they are left recursive
 * to conserve yacc stack. they need to
 * be reversed to interpret correctly
 */
xdcl_list_r:
	xdcl
|	xdcl_list_r xdcl
	{
		$$ = list($1, $2);
	}

vardcl_list_r:
	vardcl
|	vardcl_list_r ';' vardcl
	{
		$$ = nod(OLIST, $1, $3);
	}

constdcl_list_r:
	constdcl1
|	constdcl_list_r ';' constdcl1

typedcl_list_r:
	typedcl
|	typedcl_list_r ';' typedcl

structdcl_list_r:
	structdcl
	{
		$$ = cleanidlist($1);
	}
|	structdcl_list_r ';' structdcl
	{
		$$ = cleanidlist($3);
		$$ = nod(OLIST, $1, $$);
	}

interfacedcl_list_r:
	interfacedcl
	{
		$$ = cleanidlist($1);
	}
|	interfacedcl_list_r ';' interfacedcl
	{
		$$ = cleanidlist($3);
		$$ = nod(OLIST, $1, $$);
	}

structdcl:
	new_field ',' structdcl
	{
		$$ = nod(ODCLFIELD, $1, N);
		$$ = nod(OLIST, $$, $3);
	}
|	new_field type oliteral
	{
		$$ = nod(ODCLFIELD, $1, N);
		$$->type = $2;
		$$->val = $3;
	}
|	embed oliteral
	{
		$$ = $1;
		$$->val = $2;
	}
|	'*' embed oliteral
	{
		$$ = $2;
		$$->type = ptrto($$->type);
		$$->val = $3;
	}

embed:
	LATYPE
	{
		$$ = embedded($1);
	}
|	lpack '.' LATYPE
	{
		$$ = embedded($3);
		context = nil;
	}

interfacedcl1:
	new_name ',' interfacedcl1
	{
		$$ = nod(ODCLFIELD, $1, N);
		$$ = nod(OLIST, $$, $3);
	}
|	new_name indcl
	{
		$$ = nod(ODCLFIELD, $1, N);
		$$->type = $2;
	}

interfacedcl:
	interfacedcl1
|	latype
	{
		$$ = nod(ODCLFIELD, N, N);
		$$->type = oldtype($1);
	}

indcl:
	'(' oarg_type_list ')' fnres
	{
		// without func keyword
		$$ = functype(fakethis(), $2, $4);
	}

/*
 * function arguments.
 *
 * the hard part is that when we're reading a list of names,
 * we don't know if they are going to be the names of
 * parameters (like "a,b,c int") or the types of anonymous
 * parameters (like "int, string, bool").
 *
 * an arg_chunk is a comma-separated list of arguments
 * that ends in an obvious type, either "a, b, c x" or "a, b, c, *x".
 * in the first case, a, b, c are parameters of type x.
 * in the second case, a, b, c, and *x are types of anonymous parameters.
 */
arg_chunk:
	new_name_list_r type
	{
		$$ = nametodcl($1, $2);
	}
|	new_name_list_r dotdotdot
	{
		$$ = nametodcl($1, $2);
	}
|	non_name_type
	{
		$$ = anondcl($1);
	}
|	new_name_list_r ',' non_name_type
	{
		$1 = nametoanondcl($1);
		$$ = appendr($1, anondcl($3));
	}

arg_chunk_list_r:
	arg_chunk
|	arg_chunk_list_r ',' arg_chunk
	{
		$$ = appendr($1, $3);
	}

/*
 * an arg type list is a sequence of arg chunks,
 * possibly ending in a list of names (plain "a,b,c"),
 * which must be the types of anonymous parameters.
 */
arg_type_list_r:
	arg_chunk_list_r
|	arg_chunk_list_r ',' new_name_list_r
	{
		$3 = nametoanondcl($3);
		$$ = appendr($1, $3);
	}
|	new_name_list_r
	{
		$$ = nametoanondcl($1);
	}

arg_type_list:
	arg_type_list_r
	{
		$$ = rev($1);
		checkarglist($$);
	}

/*
 * statement that doesn't need semicolon terminator
 */
Astmt:
	complex_stmt
|	compound_stmt
|	Acommon_dcl
|	';'
	{
		$$ = N;
	}
|	error Astmt
	{
		$$ = N;
	}
|	labelname ':'
	{
		$$ = nod(OLABEL, $1, N);
	}
|	Bstmt ';'

/*
 * statement that does
 */
Bstmt:
	semi_stmt
|	Bcommon_dcl
|	simple_stmt

/*
 * statement list that doesn't need semicolon terminator
 */
Astmt_list_r:
	Astmt
|	Astmt_list_r Astmt
	{
		$$ = list($1, $2);
	}

/*
 * statement list that needs semicolon terminator
 */
Bstmt_list_r:
	Bstmt
|	Astmt_list_r Bstmt
	{
		$$ = list($1, $2);
	}

stmt_list_r:
	Astmt_list_r
|	Bstmt_list_r

expr_list_r:
	expr
|	expr_list_r ',' expr
	{
		$$ = nod(OLIST, $1, $3);
	}

new_name_list_r:
	new_name
|	new_name_list_r ',' new_name
	{
		$$ = nod(OLIST, $1, $3);
	}

exprsym3:
	expr
|	sym3
	{
		$$ = newname($1);
	}
|	LATYPE
	{
		$$ = newname($1);
	}

exprsym3_list_r:
	exprsym3
|	exprsym3_list_r ',' exprsym3
	{
		$$ = nod(OLIST, $1, $3);
	}

import_stmt_list_r:
	import_stmt
|	import_stmt_list_r ';' import_stmt

hidden_import_list_r:
|	hidden_import_list_r hidden_import

hidden_funarg_list_r:
	hidden_dcl
|	hidden_funarg_list_r ',' hidden_dcl
	{
		$$ = nod(OLIST, $1, $3);
	}

hidden_funarg_list:
	hidden_funarg_list_r
	{
		$$ = rev($1);
	}

hidden_structdcl_list_r:
	hidden_structdcl
|	hidden_structdcl_list_r ';' hidden_structdcl
	{
		$$ = nod(OLIST, $1, $3);
	}

hidden_structdcl_list:
	hidden_structdcl_list_r
	{
		$$ = rev($1);
	}

hidden_interfacedcl_list_r:
	hidden_interfacedcl
|	hidden_interfacedcl_list_r ';' hidden_interfacedcl
	{
		$$ = nod(OLIST, $1, $3);
	}

hidden_interfacedcl_list:
	hidden_interfacedcl_list_r
	{
		$$ = rev($1);
	}

/*
 * list of combo of keyval and val
 */
keyval_list_r:
	keyval
|	expr
|	keyval_list_r ',' keyval
	{
		$$ = nod(OLIST, $1, $3);
	}
|	keyval_list_r ',' expr
	{
		$$ = nod(OLIST, $1, $3);
	}

/*
 * have to spell this out using _r lists to avoid yacc conflict
 */
braced_keyexpr_list:
	{
		$$ = N;
	}
|	keyval_list_r ocomma
	{
		$$ = rev($1);
	}


/*
 * the one compromise of a
 * non-reversed list
 */
expr_list:
	expr_list_r
	{
		$$ = rev($1);
	}


/*
 * optional things
 */
osemi:
|	';'

ocomma:
|	','

oexpr:
	{
		$$ = N;
	}
|	expr

oexpr_list:
	{
		$$ = N;
	}
|	expr_list

osimple_stmt:
	{
		$$ = N;
	}
|	simple_stmt

ostmt_list:
	{
		$$ = N;
	}
|	stmt_list_r
	{
		$$ = rev($1);
	}

ocaseblock_list:
	{
		$$ = N;
	}
|	caseblock_list_r
	{
		$$ = rev($1);
	}

oxdcl_list:
	{
		$$ = N;
	}
|	xdcl_list_r
	{
		$$ = rev($1);
	}

oarg_type_list:
	{
		$$ = N;
	}
|	arg_type_list

ohidden_funarg_list:
	{
		$$ = N;
	}
|	hidden_funarg_list

ohidden_structdcl_list:
	{
		$$ = N;
	}
|	hidden_structdcl_list

ohidden_interfacedcl_list:
	{
		$$ = N;
	}
|	hidden_interfacedcl_list

oliteral:
	{
		$$.ctype = CTxxx;
	}
|	LLITERAL

/*
 * import syntax from header of
 * an output package
 */
hidden_import:
	LPACKAGE sym1
	/* variables */
|	LVAR hidden_pkg_importsym hidden_type
	{
		importvar($2, $3, PEXTERN);
	}
|	LCONST hidden_pkg_importsym '=' hidden_constant
	{
		importconst($2, types[TIDEAL], $4);
	}
|	LCONST hidden_pkg_importsym hidden_type '=' hidden_constant
	{
		importconst($2, $3, $5);
	}
|	LTYPE hidden_pkg_importsym hidden_type
	{
		importtype($2, $3);
	}
|	LFUNC hidden_pkg_importsym '(' ohidden_funarg_list ')' ohidden_funres
	{
		importvar($2, functype(N, $4, $6), PFUNC);
	}
|	LFUNC '(' hidden_funarg_list ')' sym1 '(' ohidden_funarg_list ')' ohidden_funres
	{
		if($3->op != ODCLFIELD) {
			yyerror("bad receiver in method");
			YYERROR;
		}
		importmethod($5, functype($3, $7, $9));
	}

hidden_type:
	hidden_type1
|	hidden_type2

hidden_type1:
	hidden_importsym
	{
		$$ = pkgtype($1->sym->name, $1->psym->name);
	}
|	LATYPE
	{
		$$ = oldtype($1);
	}
|	'[' ']' hidden_type
	{
		$$ = aindex(N, $3);
	}
|	'[' LLITERAL ']' hidden_type
	{
		$$ = aindex(nodlit($2), $4);
	}
|	LMAP '[' hidden_type ']' hidden_type
	{
		$$ = maptype($3, $5);
	}
|	LSTRUCT '{' ohidden_structdcl_list '}'
	{
		$$ = dostruct($3, TSTRUCT);
	}
|	LINTERFACE '{' ohidden_interfacedcl_list '}'
	{
		$$ = dostruct($3, TINTER);
		$$ = sortinter($$);
	}
|	'*' hidden_type
	{
		checkwidth($2);
		$$ = ptrto($2);
	}
|	LCOMM LCHAN hidden_type
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Crecv;
	}
|	LCHAN LCOMM hidden_type1
	{
		$$ = typ(TCHAN);
		$$->type = $3;
		$$->chan = Csend;
	}
|	LDDD
	{
		$$ = typ(TDDD);
	}

hidden_type2:
	LCHAN hidden_type
	{
		$$ = typ(TCHAN);
		$$->type = $2;
		$$->chan = Cboth;
	}
|	LFUNC '(' ohidden_funarg_list ')' ohidden_funres
	{
		$$ = functype(N, $3, $5);
	}

hidden_dcl:
	sym1 hidden_type
	{
		$$ = nod(ODCLFIELD, newname($1), N);
		$$->type = $2;
	}
|	'?' hidden_type
	{
		$$ = nod(ODCLFIELD, N, N);
		$$->type = $2;
	}

hidden_structdcl:
	sym1 hidden_type oliteral
	{
		$$ = nod(ODCLFIELD, newname($1), N);
		$$->type = $2;
		$$->val = $3;
	}
|	'?' hidden_type oliteral
	{
		if(isptr[$2->etype]) {
			$$ = embedded($2->type->sym);
			$$->type = ptrto($$->type);
		} else
			$$ = embedded($2->sym);
		$$->val = $3;
	}

hidden_interfacedcl:
	sym1 '(' ohidden_funarg_list ')' ohidden_funres
	{
		$$ = nod(ODCLFIELD, newname($1), N);
		$$->type = functype(fakethis(), $3, $5);
	}

ohidden_funres:
	{
		$$ = N;
	}
|	hidden_funres

hidden_funres:
	'(' ohidden_funarg_list ')'
	{
		$$ = $2;
	}
|	hidden_type1
	{
		$$ = nod(ODCLFIELD, N, N);
		$$->type = $1;
	}

hidden_constant:
	LLITERAL
	{
		$$ = nodlit($1);
	}
|	'-' LLITERAL
	{
		$$ = nodlit($2);
		switch($$->val.ctype){
		case CTINT:
			mpnegfix($$->val.u.xval);
			break;
		case CTFLT:
			mpnegflt($$->val.u.fval);
			break;
		default:
			yyerror("bad negated constant");
		}
	}
|	LTRUE
	{
		$$ = nodbool(1);
	}
|	LFALSE
	{
		$$ = nodbool(0);
	}

hidden_importsym:
	sym1 '.' sym2
	{
		$$ = nod(OIMPORT, N, N);
		$$->psym = $1;
		$$->sym = $3;
	}

hidden_pkg_importsym:
	hidden_importsym
	{
		$$ = $1;
		pkgcontext = $$->psym->name;
	}


/*
 * helpful error messages.
 * THIS SECTION MUST BE AT THE END OF THE FILE.
 *
 * these rules trigger reduce/reduce conflicts in the grammar.
 * they are safe because reduce/reduce conflicts are resolved
 * in favor of rules appearing earlier in the grammar, and these
 * are at the end of the file.
 *
 * to check whether the rest of the grammar is free of
 * reduce/reduce conflicts, comment this section out by
 * removing the slash on the next line.
 *
 * there should be exactly 1 reduce/reduce conflict
 * when this block is commented out.
 */
lpack:
	LATYPE
	{
		yyerror("%s is type, not package", $1->name);
		YYERROR;
	}

latype:
	LPACK
	{
		yyerror("%s is package, not type", $1->name);
		YYERROR;
	}
|	LNAME
	{
		yyerror("no type %s", $1->name);
		YYERROR;
	}

nametype:
	LNAME
	{
		yyerror("no type %s", $1->name);
		YYERROR;
	}

/**/

