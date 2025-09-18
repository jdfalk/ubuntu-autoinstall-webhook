<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Copilot/AI Agent Coding Instructions System](#copilotai-agent-coding-instructions-system)
  - [🚨 CRITICAL: Documentation Update Protocol](#-critical-documentation-update-protocol)
  - [System Overview](#system-overview)
  - [How It Works](#how-it-works)
  - [For Contributors](#for-contributors)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- file: .github/copilot-instructions.md -->
<!-- version: 2.0.1 -->
<!-- guid: 4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a -->

# Copilot/AI Agent Coding Instructions System

This repository uses a centralized, modular system for Copilot/AI agent coding,
documentation, and workflow instructions, following the latest VS Code Copilot
customization best practices.

## 🚨 CRITICAL: Documentation Update Protocol

This repository no longer uses doc-update scripts. Follow these rules instead:

- Edit documentation directly in the files within this repository.
- Keep the required file header (file path, version, guid) and bump the version on any change.
- Do not use create-doc-update.sh or related scripts; they are retired.
- Follow `.github/instructions/general-coding.instructions.md` and language-specific instruction files for rules.
- Prefer VS Code tasks for git operations (Git Add All, Git Commit, Git Push).

## System Overview

- **General rules**: `.github/instructions/general-coding.instructions.md`
  (applies to all files)
- **Language/task-specific rules**: `.github/instructions/*.instructions.md`
  (with `applyTo` frontmatter)
- **Prompt files**: `.github/prompts/` (for Copilot/AI prompt customization)
- **Agent-specific docs**: `.github/AGENTS.md`, `.github/CLAUDE.md`, etc.
  (pointers to this system)
- **VS Code integration**: `.vscode/copilot/` contains symlinks to canonical
  `.github/instructions/` files for VS Code Copilot features

## How It Works

- **General instructions** are always included for all files and languages.
- **Language/task-specific instructions** extend the general rules and use the
  `applyTo` field to target file globs (e.g., `**/*.go`).
- **All code style, documentation, and workflow rules are now found exclusively
  in `.github/instructions/*.instructions.md` files.**
- **Prompt files** are stored in `.github/prompts/` and can reference
  instructions as needed.
- **Agent docs** (e.g., AGENTS.md) point to `.github/` as the canonical source
  for all rules.
- **VS Code** uses symlinks in `.vscode/copilot/` to include these instructions
  for Copilot customization.

## For Contributors

- **Edit or add rules** in `.github/instructions/` only. Do not use or reference
  any `code-style-*.md` files; these are obsolete.
- **Add new prompts** in `.github/prompts/`.
- **Update agent docs** to reference this system.
- **Do not duplicate rules**; always reference the general instructions from
  specific ones.
- **See `.github/README.md`** for a human-friendly summary and contributor
  guide.

For full details, see the
[general coding instructions](instructions/general-coding.instructions.md) and
language-specific files in `.github/instructions/`.
