// https://leetcode.com/problems/add-digits/

// 258. Add Digits

// Given a non-negative integer num, repeatedly add all its digits until
// the result has only one digit.

// Example:

// Input: 38
// Output: 2
// Explanation: The process is like: 3 + 8 = 11, 1 + 1 = 2.
//              Since 2 has only one digit, return it.
// Follow up:
// Could you do it without any loop/recursion in O(1) runtime?

package main

// AddDigits():
// add digits of num. repeat the process till you are left with
// a single digit.
//
// return the single digit result
//
func AddDigits(num int) int {
	sumOnce := func(num int) int {
		n, sum := num, 0

		for n != 0 {
			sum += n % 10
			n /= 10
		}

		return sum
	}

	myNum := num
	sum := myNum

	for sum/10 != 0 {
		sum = sumOnce(myNum)
		myNum = sum
	}

	return sum
}
