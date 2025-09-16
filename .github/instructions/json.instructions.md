<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [JSON Coding Instructions](#json-coding-instructions)
  - [Core Principles](#core-principles)
  - [General Guidelines](#general-guidelines)
  - [Required File Header](#required-file-header)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- file: .github/instructions/json.instructions.md -->
<!-- version: 1.2.1 -->
<!-- guid: 3c2d1e0f-9a8b-7c6d-5e4f-3a2b1c0d9e8f -->
<!-- DO NOT EDIT: This file is managed centrally in ghcommon repository -->
<!-- To update: Create an issue/PR in jdfalk/ghcommon -->

<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
---

applyTo: "\*_/_.json"
description: |
Coding, documentation, and workflow rules for JSON files, following Google JSON style guide and general project rules. Reference this for all JSON code, documentation, and formatting in this repository. All unique content from the Google JSON Style Guide is merged here.

---
<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->

# JSON Coding Instructions

- Follow the [general coding instructions](general-coding.instructions.md).
- Follow the
  [Google JSON Style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)
  for additional best practices.
- All JSON files must begin with the required file header (see general
  instructions for details and JSON example).

## Core Principles

- Consistency: Follow the same conventions throughout all JSON files
- Readability: Format JSON to be easily readable by humans
- Validity: Ensure all JSON is well-formed and valid
- Clarity: Use meaningful property names and structure
- Simplicity: Keep JSON structure as simple as possible

## General Guidelines

- Use proper JSON syntax as defined by RFC 7159
- Use UTF-8 encoding
- Property names must be camelCase, meaningful, and descriptive
- Avoid JavaScript reserved words as property names
- Use appropriate JSON data types (string, number, boolean, array, object, null)
- Use double quotes for all strings
- Use 2 spaces for indentation
- Omit properties that don't apply or aren't available
- Write clear, concise, and valid JSON

## Required File Header

All JSONC files must begin with a standard header as described in the
[general coding instructions](general-coding.instructions.md). The **only
exception** is files with the `.json` extension (JSON without comments), which
are exempt from this requirement and do not require a file header. For standard
`.jsonc` files, include the following header:

```jsonc
// file: path/to/file.jsonc
// version: 1.0.0
// guid: 123e4567-e89b-12d3-a456-426614174000
```
