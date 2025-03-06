# Coding Style Guide

## Repository Pattern

- Use interfaces for repository definitions
- All repository methods should accept context.Context as first parameter
- Return concrete errors rather than wrapping in custom types
- Follow standard CRUD operation naming: Create, GetByID, Update, Delete, List

## Code Structure

- Place interfaces in separate files
- Use internal package for non-public code
- Group related functionality in packages
- Follow standard Go project layout

## Error Handling

- Return errors as last return value
- Use meaningful error variables
- Don't ignore errors

## Testing

- Write table-driven tests
- Use meaningful test names
- Test both success and error cases
- Mock external dependencies
