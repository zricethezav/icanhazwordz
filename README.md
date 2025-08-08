# icanhazwordz
does this string contain words?

```go
import (
    "github.com/zricethezav/icanhazwordz"
    "fmt"
)
func main() {
    searcher := icanhazwordz.NewSearcher(DefaultFilter)
    result := searcher.Find("The quick brown fox jumps over the lazy dog")
    for _, m := range result.AllMatches {
        fmt.Println(m.Word)
    }
}
```
