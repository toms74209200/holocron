## Workflow

- **Before starting work**: Run `gh atat pull` to sync TODO.md with GitHub Issues
- **When creating a PR**: Always link the related issue
  - Use `gh pr create` with issue reference in body: `Closes #<issue_number>`
  - Example: `gh pr create --title "Implement feature" --body "Closes #1"`
- **After merging PR**: The linked issue will automatically close
- **When updating TODO.md**: Run `gh atat push` to sync changes to GitHub Issues
