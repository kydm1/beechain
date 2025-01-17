// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"
	"math/big"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
)

// Stack 栈对象封装，提供常用的栈操作
type Stack struct {
	Items []*big.Int
}

// NewStack 新创建栈对象
func NewStack() *Stack {
	return &Stack{Items: make([]*big.Int, 0, params.StackLimit)}
}

// Data 返回栈中的所有底层数据
func (st *Stack) Data() []*big.Int {
	return st.Items
}

// Push 数据入栈
func (st *Stack) Push(d *big.Int) {
	st.Items = append(st.Items, d)
}

// PushN 同时压栈多个数据
func (st *Stack) PushN(ds ...*big.Int) {
	st.Items = append(st.Items, ds...)
}

// Pop 弹出栈顶数据
func (st *Stack) Pop() (ret *big.Int) {
	ret = st.Items[len(st.Items)-1]
	st.Items = st.Items[:len(st.Items)-1]
	return
}

// Len 栈长度
func (st *Stack) Len() int {
	return len(st.Items)
}

// Swap 将栈顶数据和栈中指定位置的数据互换位置
func (st *Stack) Swap(n int) {
	st.Items[st.Len()-n], st.Items[st.Len()-1] = st.Items[st.Len()-1], st.Items[st.Len()-n]
}

// Dup 复制栈中指定位置的数据的栈顶
func (st *Stack) Dup(pool *IntPool, n int) {
	st.Push(pool.Get().Set(st.Items[st.Len()-n]))
}

// Peek 返回顶端数据
func (st *Stack) Peek() *big.Int {
	return st.Items[st.Len()-1]
}

// Back 返回第n个取值
func (st *Stack) Back(n int) *big.Int {
	return st.Items[st.Len()-n-1]
}

// Require 检查栈是否满足长度要求
func (st *Stack) Require(n int) error {
	if st.Len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.Items), n)
	}
	return nil
}

// Print 印栈对象（调试用）
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.Items) > 0 {
		for i, val := range st.Items {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
