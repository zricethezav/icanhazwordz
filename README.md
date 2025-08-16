# icanhazwordz
does this string contain words?

```go
package main

import (
	"fmt"

	"github.com/zricethezav/icanhazwordz"
)

func main() {
	searcher := icanhazwordz.NewSearcher(icanhazwordz.Filter{
		MinLength: 3,
		MaxLength: 4,
	})
	result := searcher.Find("The quick brown fox jumps over the lazy dog")
	for _, m := range result.Matches {
		fmt.Println(m.Word)
	}

    fmt.Println("---")
    seracher.Filter =
	searcher := icanhazwordz.NewSearcher(icanhazwordz.Filter{
		MinLength: 3,
		MaxLength: 4,
        PreferLongestNonOverlapping: true,
	})
	result := searcher.Find("The quick brown fox jumps over the lazy dog")
	for _, m := range result.Matches {
		fmt.Println(m.Word)
	}
}
```

Output:
```
the
brow
row
own
fox
jump
ump
over
the
laz
lazy
dog
---
the
brow
fox
jump
over
the
lazy
dog
```

Shoutout to the 🐐, [github.com/BobuSumisu/aho-corasick](github.com/BobuSumisu/aho-corasick)

---
Folks, it's MIT

![a2l827](https://github.com/user-attachments/assets/6a8785d6-d34e-436a-a16b-a7be99dc7843)
