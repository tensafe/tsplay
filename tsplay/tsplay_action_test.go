package tsplay

import (
	"fmt"
	"testing"
)

func Test_navigate(L *testing.T) {

}

func Test_get_text(L *testing.T) {
	script_single := `
-- 获取单个元素的文本内容
local singleText = get_text("h1")
if singleText then
    print("Single Text:", singleText)
end
`
	script_table := `
-- 获取多个元素的文本内容
local multipleTexts = get_text("p")
if type(multipleTexts) == "table" then
    for i, text in ipairs(multipleTexts) do
        print("Text", i, ":", text)
    end
end
`
	fmt.Printf(script_single, script_table)
}
