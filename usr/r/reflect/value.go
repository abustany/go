// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Reflection library.
// Handling values.

package reflect

import (
	"reflect";
)


type Addr uint64	// TODO: where are ptrint/intptr etc?

export type Value interface {
	Kind()	int;
	Type()	Type;
}

func NewValueAddr(typ Type, addr Addr) Value

type Creator *(typ Type, addr Addr) Value

// Conversion functions, implemented in assembler
func AddrToPtrAddr(Addr) *Addr
func AddrToPtrInt8(Addr) *int8
func AddrToPtrInt16(Addr) *int16
func AddrToPtrInt32(Addr) *int32
func AddrToPtrInt64(Addr) *int64
func AddrToPtrUint8(Addr) *uint8
func PtrUint8ToAddr(*uint8) Addr
func AddrToPtrUint16(Addr) *uint16
func AddrToPtrUint32(Addr) *uint32
func AddrToPtrUint64(Addr) *uint64
func PtrUint64ToAddr(*uint64) Addr
func AddrToPtrFloat32(Addr) *float32
func AddrToPtrFloat64(Addr) *float64
func AddrToPtrFloat80(Addr) *float80
func AddrToPtrString(Addr) *string

// -- Int8

export type Int8Value interface {
	Kind()	int;
	Get()	int8;
	Put(int8);
	Type()	Type;
}

type Int8ValueStruct struct {
	addr	Addr
}

func (v *Int8ValueStruct) Kind() int {
	return Int8Kind
}

func (v *Int8ValueStruct) Type() Type {
	return Int8
}

func (v *Int8ValueStruct) Get() int8 {
	return *AddrToPtrInt8(v.addr)
}

func (v *Int8ValueStruct) Put(i int8) {
	*AddrToPtrInt8(v.addr) = i
}

func Int8Creator(typ Type, addr Addr) Value {
	v := new(Int8ValueStruct);
	v.addr = addr;
	return v;
}

// -- Int16

export type Int16Value interface {
	Kind()	int;
	Get()	int16;
	Put(int16);
	Type()	Type;
}

type Int16ValueStruct struct {
	addr	Addr
}

func (v *Int16ValueStruct) Kind() int {
	return Int16Kind
}

func (v *Int16ValueStruct) Type() Type {
	return Int16
}

func (v *Int16ValueStruct) Get() int16 {
	return *AddrToPtrInt16(v.addr)
}

func (v *Int16ValueStruct) Put(i int16) {
	*AddrToPtrInt16(v.addr) = i
}

func Int16Creator(typ Type, addr Addr) Value {
	v := new(Int16ValueStruct);
	v.addr = addr;
	return v;
}

// -- Int32

export type Int32Value interface {
	Kind()	int;
	Get()	int32;
	Put(int32);
	Type()	Type;
}

type Int32ValueStruct struct {
	addr	Addr
}

func (v *Int32ValueStruct) Type() Type {
	return Int32
}

func (v *Int32ValueStruct) Kind() int {
	return Int32Kind
}

func (v *Int32ValueStruct) Get() int32 {
	return *AddrToPtrInt32(v.addr)
}

func (v *Int32ValueStruct) Put(i int32) {
	*AddrToPtrInt32(v.addr) = i
}

func Int32Creator(typ Type, addr Addr) Value {
	v := new(Int32ValueStruct);
	v.addr = addr;
	return v;
}

// -- Int64

export type Int64Value interface {
	Kind()	int;
	Get()	int64;
	Put(int64);
	Type()	Type;
}

type Int64ValueStruct struct {
	addr	Addr
}

func (v *Int64ValueStruct) Kind() int {
	return Int64Kind
}

func (v *Int64ValueStruct) Type() Type {
	return Int64
}

func (v *Int64ValueStruct) Get() int64 {
	return *AddrToPtrInt64(v.addr)
}

func (v *Int64ValueStruct) Put(i int64) {
	*AddrToPtrInt64(v.addr) = i
}

func Int64Creator(typ Type, addr Addr) Value {
	v := new(Int64ValueStruct);
	v.addr = addr;
	return v;
}

// -- Uint8

export type Uint8Value interface {
	Kind()	int;
	Get()	uint8;
	Put(uint8);
	Type()	Type;
}

type Uint8ValueStruct struct {
	addr	Addr
}

func (v *Uint8ValueStruct) Kind() int {
	return Uint8Kind
}

func (v *Uint8ValueStruct) Type() Type {
	return Uint8
}

func (v *Uint8ValueStruct) Get() uint8 {
	return *AddrToPtrUint8(v.addr)
}

func (v *Uint8ValueStruct) Put(i uint8) {
	*AddrToPtrUint8(v.addr) = i
}

func Uint8Creator(typ Type, addr Addr) Value {
	v := new(Uint8ValueStruct);
	v.addr = addr;
	return v;
}

// -- Uint16

export type Uint16Value interface {
	Kind()	int;
	Get()	uint16;
	Put(uint16);
	Type()	Type;
}

type Uint16ValueStruct struct {
	addr	Addr
}

func (v *Uint16ValueStruct) Kind() int {
	return Uint16Kind
}

func (v *Uint16ValueStruct) Type() Type {
	return Uint16
}

func (v *Uint16ValueStruct) Get() uint16 {
	return *AddrToPtrUint16(v.addr)
}

func (v *Uint16ValueStruct) Put(i uint16) {
	*AddrToPtrUint16(v.addr) = i
}

func Uint16Creator(typ Type, addr Addr) Value {
	v := new(Uint16ValueStruct);
	v.addr = addr;
	return v;
}

// -- Uint32

export type Uint32Value interface {
	Kind()	int;
	Get()	uint32;
	Put(uint32);
	Type()	Type;
}

type Uint32ValueStruct struct {
	addr	Addr
}

func (v *Uint32ValueStruct) Kind() int {
	return Uint32Kind
}

func (v *Uint32ValueStruct) Type() Type {
	return Uint32
}

func (v *Uint32ValueStruct) Get() uint32 {
	return *AddrToPtrUint32(v.addr)
}

func (v *Uint32ValueStruct) Put(i uint32) {
	*AddrToPtrUint32(v.addr) = i
}

func Uint32Creator(typ Type, addr Addr) Value {
	v := new(Uint32ValueStruct);
	v.addr = addr;
	return v;
}

// -- Uint64

export type Uint64Value interface {
	Kind()	int;
	Get()	uint64;
	Put(uint64);
	Type()	Type;
}

type Uint64ValueStruct struct {
	addr	Addr
}

func (v *Uint64ValueStruct) Kind() int {
	return Uint64Kind
}

func (v *Uint64ValueStruct) Type() Type {
	return Uint64
}

func (v *Uint64ValueStruct) Get() uint64 {
	return *AddrToPtrUint64(v.addr)
}

func (v *Uint64ValueStruct) Put(i uint64) {
	*AddrToPtrUint64(v.addr) = i
}

func Uint64Creator(typ Type, addr Addr) Value {
	v := new(Uint64ValueStruct);
	v.addr = addr;
	return v;
}

// -- Float32

export type Float32Value interface {
	Kind()	int;
	Get()	float32;
	Put(float32);
	Type()	Type;
}

type Float32ValueStruct struct {
	addr	Addr
}

func (v *Float32ValueStruct) Kind() int {
	return Float32Kind
}

func (v *Float32ValueStruct) Type() Type {
	return Float32
}

func (v *Float32ValueStruct) Get() float32 {
	return *AddrToPtrFloat32(v.addr)
}

func (v *Float32ValueStruct) Put(f float32) {
	*AddrToPtrFloat32(v.addr) = f
}

func Float32Creator(typ Type, addr Addr) Value {
	v := new(Float32ValueStruct);
	v.addr = addr;
	return v;
}

// -- Float64

export type Float64Value interface {
	Kind()	int;
	Get()	float64;
	Put(float64);
	Type()	Type;
}

type Float64ValueStruct struct {
	addr	Addr
}

func (v *Float64ValueStruct) Kind() int {
	return Float64Kind
}

func (v *Float64ValueStruct) Type() Type {
	return Float64
}

func (v *Float64ValueStruct) Get() float64 {
	return *AddrToPtrFloat64(v.addr)
}

func (v *Float64ValueStruct) Put(f float64) {
	*AddrToPtrFloat64(v.addr) = f
}

func Float64Creator(typ Type, addr Addr) Value {
	v := new(Float64ValueStruct);
	v.addr = addr;
	return v;
}

// -- Float80

export type Float80Value interface {
	Kind()	int;
	Get()	float80;
	Put(float80);
	Type()	Type;
}

type Float80ValueStruct struct {
	addr	Addr
}

func (v *Float80ValueStruct) Kind() int {
	return Float80Kind
}

func (v *Float80ValueStruct) Type() Type {
	return Float80
}

/*
BUG: can't gen code for float80s
func (v *Float80ValueStruct) Get() float80 {
	return *AddrToPtrFloat80(v.addr)
	return 0;
}

func (v *Float80ValueStruct) Put(f float80) {
	*AddrToPtrFloat80(v.addr) = f
}
*/

func Float80Creator(typ Type, addr Addr) Value {
	v := new(Float80ValueStruct);
	v.addr = addr;
	return v;
}

// -- String

export type StringValue interface {
	Kind()	int;
	Get()	string;
	Put(string);
	Type()	Type;
}

type StringValueStruct struct {
	addr	Addr
}

func (v *StringValueStruct) Kind() int {
	return StringKind
}

func (v *StringValueStruct) Type() Type {
	return String
}

func (v *StringValueStruct) Get() string {
	return *AddrToPtrString(v.addr)
}

func (v *StringValueStruct) Put(s string) {
	*AddrToPtrString(v.addr) = s
}

func StringCreator(typ Type, addr Addr) Value {
	v := new(StringValueStruct);
	v.addr = addr;
	return v;
}

// -- Pointer

export type PtrValue interface {
	Kind()	int;
	Sub()	Value;
	Type()	Type;
	Indirect()	Addr;
}

type PtrValueStruct struct {
	addr	Addr;
	typ	Type;
}

func (v *PtrValueStruct) Kind() int {
	return PtrKind
}

func (v *PtrValueStruct) Type() Type {
	return v.typ
}

func (v *PtrValueStruct) Indirect() Addr {
	return *AddrToPtrAddr(v.addr)
}

func (v *PtrValueStruct) Sub() Value {
	return NewValueAddr(v.typ.(PtrType).Sub(), v.Indirect());
}

func PtrCreator(typ Type, addr Addr) Value {
	return &PtrValueStruct{addr, typ};
}

// -- Array	TODO: finish and test

export type ArrayValue interface {
	Kind()	int;
	Type()	Type;
	Open()	bool;
	Len()	uint64;
	Elem(i uint64)	Value;
}

type OpenArrayValueStruct struct {
	addr	Addr;
	typ	Type;
	elemtype	Type;
	elemsize	uint64;
}
/*
	Run-time representation of open arrays looks like this:
		struct	Array {
			byte*	array;		// actual data
			uint32	nel;		// number of elements
		};
*/

func (v *OpenArrayValueStruct) Kind() int {
	return ArrayKind
}

func (v *OpenArrayValueStruct) Type() Type {
	return v.typ
}

func (v *OpenArrayValueStruct) Open() bool {
	return true
}

func (v *OpenArrayValueStruct) Len() uint64 {
	return uint64(*AddrToPtrInt32(v.addr+8));
}

func (v *OpenArrayValueStruct) Elem(i uint64) Value {
	base := *AddrToPtrAddr(v.addr);
	return NewValueAddr(v.elemtype, base + i * v.elemsize);
}

type FixedArrayValueStruct struct {
	addr	Addr;
	typ	Type;
	elemtype	Type;
	elemsize	uint64;
	len	uint64;
}

func (v *FixedArrayValueStruct) Kind() int {
	return ArrayKind
}

func (v *FixedArrayValueStruct) Type() Type {
	return v.typ
}

func (v *FixedArrayValueStruct) Open() bool {
	return false
}

func (v *FixedArrayValueStruct) Len() uint64 {
	return v.len
}

func (v *FixedArrayValueStruct) Elem(i uint64) Value {
	return NewValueAddr(v.elemtype, v.addr + i * v.elemsize);
	return nil
}

func ArrayCreator(typ Type, addr Addr) Value {
	arraytype := typ.(ArrayType);
	if arraytype.Open() {
		v := new(OpenArrayValueStruct);
		v.addr = addr;
		v.typ = typ;
		v.elemtype = arraytype.Elem();
		v.elemsize = v.elemtype.Size();
		return v;
	}
	v := new(FixedArrayValueStruct);
	v.addr = addr;
	v.typ = typ;
	v.elemtype = arraytype.Elem();
	v.elemsize = v.elemtype.Size();
	v.len = arraytype.Len();
	return v;
}

// -- Map	TODO: finish and test

export type MapValue interface {
	Kind()	int;
	Type()	Type;
	Len()	int;
	Elem(key Value)	Value;
}

type MapValueStruct struct {
	addr	Addr;
	typ	Type;
	len	int;
}

func (v *MapValueStruct) Kind() int {
	return MapKind
}

func (v *MapValueStruct) Type() Type {
	return v.typ
}

func (v *MapValueStruct) Len() int {
	return v.len	// TODO: probably want this to be dynamic
}

func (v *MapValueStruct) Elem(key Value) Value {
	panic("map value element");
	return nil
}

func MapCreator(typ Type, addr Addr) Value {
	arraytype := typ.(MapType);
	v := new(MapValueStruct);
	v.addr = addr;
	v.typ = typ;
	return v;
}

// -- Chan

export type ChanValue interface {
	Kind()	int;
	Type()	Type;
}

type ChanValueStruct struct {
	addr	Addr;
	typ	Type;
	len	int;
}

func (v *ChanValueStruct) Kind() int {
	return ChanKind
}

func (v *ChanValueStruct) Type() Type {
	return v.typ
}

func ChanCreator(typ Type, addr Addr) Value {
	v := new(ChanValueStruct);
	v.addr = addr;
	v.typ = typ;
	return v;
}

// -- Struct

export type StructValue interface {
	Kind()	int;
	Type()	Type;
	Len()	int;
	Field(i int)	Value;
}

type StructValueStruct struct {
	addr	Addr;
	typ	Type;
	field	*[]Value;
}

func (v *StructValueStruct) Kind() int {
	return StructKind
}

func (v *StructValueStruct) Type() Type {
	return v.typ
}

func (v *StructValueStruct) Len() int {
	return len(v.field)
}

func (v *StructValueStruct) Field(i int) Value {
	return v.field[i]
}

func StructCreator(typ Type, addr Addr) Value {
	t := typ.(StructType);
	v := new(StructValueStruct);
	v.addr = addr;
	nfield := t.Len();
	v.field = new([]Value, nfield);
	for i := 0; i < nfield; i++ {
		name, ftype, offset := t.Field(i);
		v.field[i] = NewValueAddr(ftype, addr + offset);
	}
	v.typ = typ;
	return v;
}

// -- Interface

export type InterfaceValue interface {
	Kind()	int;
	Type()	Type;
}

type InterfaceValueInterface struct {
	addr	Addr;
	typ	Type;
}

func (v *InterfaceValueInterface) Kind() int {
	return InterfaceKind
}

func (v *InterfaceValueInterface) Type() Type {
	return v.typ
}

func InterfaceCreator(typ Type, addr Addr) Value {
	v := new(InterfaceValueInterface);
	v.addr = addr;
	v.typ = typ;
	return v;
}

// -- Func

export type FuncValue interface {
	Kind()	int;
	Type()	Type;
}

type FuncValueFunc struct {
	addr	Addr;
	typ	Type;
}

func (v *FuncValueFunc) Kind() int {
	return FuncKind
}

func (v *FuncValueFunc) Type() Type {
	return v.typ
}

func FuncCreator(typ Type, addr Addr) Value {
	v := new(FuncValueFunc);
	v.addr = addr;
	v.typ = typ;
	return v;
}

var creator *map[int] Creator

func init() {
	creator = new(map[int] Creator);
	creator[Int8Kind] = &Int8Creator;
	creator[Int16Kind] = &Int16Creator;
	creator[Int32Kind] = &Int32Creator;
	creator[Int64Kind] = &Int64Creator;
	creator[Uint8Kind] = &Uint8Creator;
	creator[Uint16Kind] = &Uint16Creator;
	creator[Uint32Kind] = &Uint32Creator;
	creator[Uint64Kind] = &Uint64Creator;
	creator[Float32Kind] = &Float32Creator;
	creator[Float64Kind] = &Float64Creator;
	creator[Float80Kind] = &Float80Creator;
	creator[StringKind] = &StringCreator;
	creator[PtrKind] = &PtrCreator;
	creator[ArrayKind] = &ArrayCreator;
	creator[MapKind] = &MapCreator;
	creator[ChanKind] = &ChanCreator;
	creator[StructKind] = &StructCreator;
	creator[InterfaceKind] = &InterfaceCreator;
	creator[FuncKind] = &FuncCreator;
}

func NewValueAddr(typ Type, addr Addr) Value {
	c, ok := creator[typ.Kind()];
	if !ok {
		panicln("no creator for type" , typ.Kind());
	}
	return c(typ, addr);
}

export func NewInitValue(typ Type) Value {
	// Some values cannot be made this way.
	switch typ.Kind() {
	case FuncKind, ChanKind, MapKind:	// must be pointers, at least for now (TODO?)
		return nil;
	case ArrayKind:
		if typ.(ArrayType).Open() {
			return nil
		}
	}
	size := typ.Size();
	if size == 0 {
		size = 1;
	}
	data := new([]uint8, size);
	return NewValueAddr(typ, PtrUint8ToAddr(&data[0]));
}

export type Empty interface {}

export func NewValue(e Empty) Value {
	value, typestring  := sys.reflect(e);
	typ := ParseTypeString("", typestring);
	// Content of interface is a value; need a permanent copy to take its address
	// so we can modify the contents. Values contain pointers to 'values'.
	ap := new(uint64);
	*ap = value;
	return NewValueAddr(typ, PtrUint64ToAddr(ap));
}
