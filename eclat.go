package main

import (
	"bufio"
	"fmt"
	"github.com/deckarep/golang-set"
	"os"
	"sort"
	"strconv"
	"strings"
)

const window = 10

type Itemset map[int][]int

type Pair struct {
	Key, Num int
}
type Pairlist []Pair

func (a Pairlist) Len() int {
	return len(a)
}
func (a Pairlist) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Pairlist) Less(i, j int) bool {
	return a[i].Num < a[j].Num
}

func getCommonSet(left []int, right []int) []int {
	lset := mapset.NewSet()
	rset := mapset.NewSet()
	for _, v := range left {
		lset.Add(v)
	}
	for _, v := range right {
		rset.Add(v)
	}
	lset = lset.Intersect(rset)
	list := make([]int, 0, lset.Cardinality())
	for v := range lset.Iter() {
		if num, ok := v.(int); ok {
			list = append(list, num)
		} else {
			return []int{}
		}
	}
	return list
}
func Eclat(minsup int, prefix []int, is *Itemset) error {
	var err error
	list := make(Pairlist, 0)
	for key, value := range *is {
		list = append(list, Pair{key, len(value)})
	}
	sort.Sort(sort.Reverse(list))
	for i, p := range list {
		isupp := len((*is)[p.Key])
		if isupp < minsup {
			break
		}
		fmt.Printf("%d:", minsup)
		fmt.Println(append(prefix, p.Key))
		suffix := make(Itemset, 0)
		for _, pp := range list[i+1 : len(list)] {
			comlist := getCommonSet((*is)[p.Key], (*is)[pp.Key])
			if len(comlist) >= minsup {
				suffix[pp.Key] = comlist
			}
		}
		Eclat(minsup, append(prefix, p.Key), &suffix)
	}
	return err
}

func scanTransaction(ifile string, is *Itemset) error {
	var fp *os.File
	var err error
	if len(ifile) == 0 {
		return err
	}
	if fp, err = os.Open(ifile); err != nil {
		return err
	}
	scanner := bufio.NewScanner(fp)
	for trans := 0; scanner.Scan(); trans++ {
		s := strings.Split(scanner.Text(), " ")
		if len(s) == 0 {
			break
		}
		for _, element := range s {
			item, _ := strconv.Atoi(element)
			if _, ok := (*is)[item]; !ok {
				(*is)[item] = make([]int, 0, window)
			}
			(*is)[item] = append((*is)[item], trans)
		}
	}
	fmt.Println(*is)
	return scanner.Err()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("No arguments.")
		return
	}
	num, _ := strconv.Atoi(os.Args[2])
	if num < 1 {
		fmt.Printf("Minimum support error (%d).\n", num)
		return
	}
	is := make(Itemset)
	var prefix []int
	if err := scanTransaction(os.Args[1], &is); err != nil {
		fmt.Println(err)
		return
	}
	if err := Eclat(num, prefix, &is); err != nil {
		fmt.Println(err)
	}
}
