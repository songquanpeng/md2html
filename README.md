# Markdown to HTML

## TODO
- [x] Multi-level list.
- [ ] Add style for code block.
- [ ] Support code block without triple backticks.

## Grammar
```
article -> section_list
section_list -> section
              | section + section_list
section -> title
         | dividing_line
         | content
         | quote
         | code_block
         | uncompleted_task_list
         | completed_task_list
         | unordered_list
         | ordered_list
title ->  TitleToken + content
dividing_line -> DividingLineToken
content -> TextToken + rich_text + TextToken
rich_text -> italic
           | bold
           | code
           | strikethrough
           | link
           | image
italic -> SingleStarToken + content + SingleStarToken
        | SingleUnderscoreToken + content + SingleUnderscoreToken
bold -> DoubleStarToken + content + DoubleStarToken
      | DoubleUnderscoreToken + content + DoubleUnderscoreToken
code -> SingleBacktickToken + TextToken + SingleBacktickToken
strikethrough -> DoubleTildeToken + content + DoubleTildeToken
link -> LinkHeadToken + content + LinkBodyToken
image -> ImageHeadToken + text + LinkBodyToken
quote -> QuoteToken + content
code_block -> CodeBlockToken
uncompleted_task_list -> UncompletedTaskToken + content
completed_task_list -> CompletedTaskToken + content
unordered_list -> UnorderedListToken + content
ordered_list -> OrderedListToken + content
```