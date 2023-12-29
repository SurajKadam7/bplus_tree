# B+ Tree Implementation in Go Generics

## Overview

This repository contains an implementation of a B+ tree using Go generics. The B+ tree is a powerful data structure for indexing and organizing data, and this project aims to provide a generic and efficient implementation.

## Features

- **Go Generics Implementation**: Utilizes Go's generics to create a versatile and type-safe B+ tree.
- **Documentation**: Well-documented code with inline comments explaining the logic and functionalities.
- **Test Cases**: Comprehensive test suite to ensure correctness and robustness.
- **Randomized Tests**: Incorporates randomized tests to cover various scenarios and edge cases.

## Usage

### Interface (BPTree)

The project exposes an interface, `BPTree`, for clients to interact with the B+ tree. Here's an overview of the methods available in the interface:

```go
type Comparable interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int32 | int64 | float32 | float64
}

type BPTree[T1 Comparable, T2 any] interface {
	Get(key T1) (value T2, err error)
	Put(key T1, value T2) (err error)
	Delete(key T1) (err error)
}
```


### Example Usage
Here's an example demonstrating how to use the B+ tree:

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	t := New[int, float64](3)

	err := t.Put(10, 20.0)
	if err != nil {
		log.Println(err)
	}

	value, err := t.Get(10)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("key : %d value : %f \n", 10, value)

	err = t.Delete(10)
	if err != nil {
		log.Println(err)
	}
}
```

# Contribution Guidelines

Contributions are welcome! If you find issues, have suggestions for improvements, or want to add new features, feel free to create a pull request. Please follow these guidelines:

1. **Fork the repository.**
2. **Create a new branch for your feature/fix.**
3. **Make changes and add test cases.**
4. **Create a pull request with a clear description of changes.**

Thank you for considering contributing to this project!

# About Me

- **Email:** [Email](surajskadam7@gmail.com)
- **LinkedIn:** [LinkedIn](https://www.linkedin.com/in/suraj-kadam-4b549a208/)

