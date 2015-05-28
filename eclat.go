// This implementation is following an eclat python script found here http://adrem.ua.ac.be/~goethals/software/.
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

// window is the initial capacity of transaction list in Itemset.
const window = 10

// Itemset is a map of int keys for item and []int values for transaction lists.
type Itemset map[int][]int

type pair struct {
	Key, Num int
}
type pairlist []pair

func (a pairlist) Len() int {
	return len(a)
}
func (a pairlist) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a pairlist) Less(i, j int) bool {
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

func printSet(isup int, prefix []int) {
	sort.Ints(prefix)
	fmt.Printf("[")
	for i, m := range prefix {
		if i == len(prefix)-1 {
			fmt.Printf("'%d'] :", m)
		} else {
			fmt.Printf("'%d', ", m)
		}
	}
	fmt.Printf(" %d\n", isup)
}

// Eclat outputs all itemsets whose support number is equal or greater than minsup.
// Format: ['item1', 'item2', ...] : support number
func Eclat(minsup int, prefix []int, is *Itemset) error {
	var err error
	list := make(pairlist, 0)
	for key, value := range *is {
		list = append(list, pair{key, len(value)})
	}
	sort.Sort(sort.Reverse(list))
	for i, p := range list {
		isup := len((*is)[p.Key])
		if isup < minsup {
			break
		}
		printSet(isup, append(prefix, p.Key))
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
