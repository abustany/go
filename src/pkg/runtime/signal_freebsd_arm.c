// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "defs_GOOS_GOARCH.h"
#include "signals_GOOS.h"
#include "os_GOOS.h"

void
runtime·dumpregs(Mcontext *r)
{
	runtime·printf("r0      %x\n", r->r0);
	runtime·printf("r1      %x\n", r->r1);
	runtime·printf("r2      %x\n", r->r2);
	runtime·printf("r3      %x\n", r->r3);
	runtime·printf("r4      %x\n", r->r4);
	runtime·printf("r5      %x\n", r->r5);
	runtime·printf("r6      %x\n", r->r6);
	runtime·printf("r7      %x\n", r->r7);
	runtime·printf("r8      %x\n", r->r8);
	runtime·printf("r9      %x\n", r->r9);
	runtime·printf("r10     %x\n", r->r10);
	runtime·printf("fp      %x\n", r->r11);
	runtime·printf("ip      %x\n", r->r12);
	runtime·printf("sp      %x\n", r->r13);
	runtime·printf("lr      %x\n", r->r14);
	runtime·printf("pc      %x\n", r->r15);
	runtime·printf("cpsr    %x\n", r->cpsr);
}

extern void runtime·sigtramp(void);

typedef struct sigaction {
	union {
		void    (*__sa_handler)(int32);
		void    (*__sa_sigaction)(int32, Siginfo*, void *);
	} __sigaction_u;		/* signal handler */
	int32	sa_flags;		/* see signal options below */
	int64	sa_mask;		/* signal mask to apply */
} Sigaction;

void
runtime·sighandler(int32 sig, Siginfo *info, void *context, G *gp)
{
	Ucontext *uc;
	Mcontext *r;
	SigTab *t;

	uc = context;
	r = &uc->uc_mcontext;

	if(sig == SIGPROF) {
		runtime·sigprof((uint8*)r->r15, (uint8*)r->r13, (uint8*)r->r14, gp);
		return;
	}

	t = &runtime·sigtab[sig];
	if(info->si_code != SI_USER && (t->flags & SigPanic)) {
		if(gp == nil)
			goto Throw;
		// Make it look like a call to the signal func.
		// Have to pass arguments out of band since
		// augmenting the stack frame would break
		// the unwinding code.
		gp->sig = sig;
		gp->sigcode0 = info->si_code;
		gp->sigcode1 = (uintptr)info->si_addr;
		gp->sigpc = r->r15;

		// Only push runtime·sigpanic if r->mc_rip != 0.
		// If r->mc_rip == 0, probably panicked because of a
		// call to a nil func.  Not pushing that onto sp will
		// make the trace look like a call to runtime·sigpanic instead.
		// (Otherwise the trace will end at runtime·sigpanic and we
		// won't get to see who faulted.)
		if(r->r15 != 0)
			r->r14 = r->r15;
		// In case we are panicking from external C code
		r->r10 = (uintptr)gp;
		r->r9 = (uintptr)m;
		r->r15 = (uintptr)runtime·sigpanic;
		return;
	}

	if(info->si_code == SI_USER || (t->flags & SigNotify))
		if(runtime·sigsend(sig))
			return;
	if(t->flags & SigKill)
		runtime·exit(2);
	if(!(t->flags & SigThrow))
		return;

Throw:
	runtime·startpanic();

	if(sig < 0 || sig >= NSIG)
		runtime·printf("Signal %d\n", sig);
	else
		runtime·printf("%s\n", runtime·sigtab[sig].name);

	runtime·printf("PC=%x\n", r->r15);
	runtime·printf("\n");

	if(runtime·gotraceback()){
		runtime·traceback((void*)r->r15, (void*)r->r13, (void*)r->r14, gp);
		runtime·tracebackothers(gp);
		runtime·printf("\n");
		runtime·dumpregs(r);
	}

//	breakpoint();
	runtime·exit(2);
}

void
runtime·signalstack(byte *p, int32 n)
{
	Sigaltstack st;

	st.ss_sp = (int8*)p;
	st.ss_size = n;
	st.ss_flags = 0;
	runtime·sigaltstack(&st, nil);
}

void
runtime·setsig(int32 i, void (*fn)(int32, Siginfo*, void*, G*), bool restart)
{
	Sigaction sa;

	runtime·memclr((byte*)&sa, sizeof sa);
	sa.sa_flags = SA_SIGINFO|SA_ONSTACK;
	if(restart)
		sa.sa_flags |= SA_RESTART;
	sa.sa_mask = ~0ULL;
	if (fn == runtime·sighandler)
		fn = (void*)runtime·sigtramp;
	sa.__sigaction_u.__sa_sigaction = (void*)fn;
	runtime·sigaction(i, &sa, nil);
}

void
runtime·checkgoarm(void)
{
	// TODO(minux)
}

#pragma textflag 7
int64
runtime·cputicks() {
	// Currently cputicks() is used in blocking profiler and to seed runtime·fastrand1().
	// runtime·nanotime() is a poor approximation of CPU ticks that is enough for the profiler.
	// TODO: need more entropy to better seed fastrand1.
	return runtime·nanotime();
}
