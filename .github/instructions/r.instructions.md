<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [R Coding Instructions](#r-coding-instructions)
  - [Naming Conventions](#naming-conventions)
  - [Syntax and Formatting](#syntax-and-formatting)
  - [Required File Header](#required-file-header)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- file: .github/instructions/r.instructions.md -->
<!-- version: 1.2.1 -->
<!-- guid: 6c5b4a3c-2d1e-0f9a-8b7c-6d5e4f3a2b1c -->
<!-- DO NOT EDIT: This file is managed centrally in ghcommon repository -->
<!-- To update: Create an issue/PR in jdfalk/ghcommon -->

<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
---

applyTo: "\*_/_.R"
description: |
Coding, documentation, and workflow rules for R files, following Google/Tidyverse R style guide and general project rules. Reference this for all R code, documentation, and formatting in this repository. All unique content from the Google/Tidyverse R Style Guide is merged here.

---
<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->

# R Coding Instructions

- Follow the [general coding instructions](general-coding.instructions.md).
- Follow the
  [Google R Style Guide](https://google.github.io/styleguide/Rguide.html) and
  [Tidyverse Style Guide](https://style.tidyverse.org/) for additional best
  practices.
- All R files must begin with the required file header (see general instructions
  for details and R example).

## Naming Conventions

- Use snake_case for variables and functions
- Variable names should be nouns, function names should be verbs
- Google prefers BigCamelCase for functions, Tidyverse prefers snake_case
- Private functions should start with a dot
- Avoid problematic names and reserved words

## Syntax and Formatting

- Always space after commas, never before
- No spaces inside parentheses for function calls
- Most infix operators surrounded by spaces
- Limit lines to 80 characters
- Break long function calls with one argument per line
- Use 2 spaces for indentation
- Use braced expressions for control flow and loops
- Use `&&` and `||` in conditions
- Use descriptive names, avoid single character names

## Required File Header

All R files must begin with a standard header as described in the
[general coding instructions](general-coding.instructions.md). Example for R:

```r
# file: path/to/file.R
# version: 1.0.0
# guid: 123e4567-e89b-12d3-a456-426614174000
```
