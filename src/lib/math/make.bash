# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#!/bin/bash

# should probably just be a makefile

bash g1 && cp math.a $GOROOT/pkg/math.a
