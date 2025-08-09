# Command Aliases

This document lists all available command aliases for easier usage.

## Get Commands

### Index Templates
- **Primary Command**: `searchctl get index-templates`
- **Aliases**: 
  - `idx-templates`
  - `template`
  - `it`
  - `index-template`
  - `indextemplates`
  - `indextemplate`

**Examples:**
```bash
searchctl get index-templates
searchctl get idx-templates
searchctl get it
searchctl get template logs-*
```

### Component Templates
- **Primary Command**: `searchctl get component-templates`
- **Aliases**:
  - `componenttemplates`
  - `component-template`
  - `componenttemplate`
  - `ct`
  - `comp-templates`
  - `comp-template`

**Examples:**
```bash
searchctl get component-templates
searchctl get ct
searchctl get comp-template
searchctl get component-template base-*
```

## Delete Commands

### Index Template
- **Primary Command**: `searchctl delete index-template`
- **Aliases**:
  - `template`
  - `it`
  - `index-templates`
  - `indextemplates`
  - `indextemplate`

**Examples:**
```bash
searchctl delete index-template logs-template
searchctl delete template logs-template
searchctl delete it logs-template
```

### Component Template
- **Primary Command**: `searchctl delete component-template`
- **Aliases**:
  - `componenttemplate`
  - `ct`
  - `component-templates`
  - `componenttemplates`
  - `comp-template`
  - `comp-templates`

**Examples:**
```bash
searchctl delete component-template base-settings -y
searchctl delete ct base-settings -y
searchctl delete comp-template base-settings -y
```

## Usage Tips

1. **Short Aliases**: Use `it` for index templates and `ct` for component templates for quick operations
2. **Clear Naming**: Use `idx-templates` when you want to be explicit about index templates
3. **Flexibility**: All aliases work with the same parameters as the primary commands
4. **Consistency**: Aliases work across all output formats (`-o json`, `-o yaml`, etc.)
5. **Tab Completion**: Most shells support tab completion with these aliases

## Examples by Use Case

### Quick Listing
```bash
searchctl get it        # List all index templates
searchctl get ct        # List all component templates
```

### Pattern Matching
```bash
searchctl get idx-templates logs-*     # Get index templates matching pattern
searchctl get ct base-*                # Get component templates matching pattern
```

### Different Output Formats
```bash
searchctl get it -o json           # JSON output
searchctl get ct -o yaml           # YAML output
```

### Quick Deletion
```bash
searchctl delete it old-template -y      # Delete index template
searchctl delete ct old-component -y     # Delete component template
```