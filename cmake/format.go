// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package cmake

import "fmt"

func Quote(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
