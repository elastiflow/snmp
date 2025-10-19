# snmp

Definitions for the ElastiFlow&trade; SNMP Poller and SNMP Trap Collectors.

> NOTE: The `main` branch requires `7.18.0` or later of the ElastiFlow NetObserv Collectors.

## Validation

The `make validate` command can be used to validate the SNMP definition files in the repository. It checks:

- SNMP Trap Enterprises
- SNMP Definitions (devices, device groups, object groups, and objects)
- SNMP Enums (integer, bitmap, and OID enums)

### When to run validations

Run `make validate` in the following situations:

- After adding new SNMP definition files
- After modifying existing SNMP definition files
- Before submitting a pull request
- As part of your CI/CD pipeline

This ensures that all SNMP definitions are correctly formatted and consistent with each other, preventing potential issues when the definitions are used by the NetObserv SNMP and SNMP Trap Collectors.

### How to run validations

Simply run the following command from the root of the repository:

```bash
make validate
```

If validation succeeds, you'll see output indicating the number of validated items for each category. If validation fails, error messages will be displayed indicating the issues that need to be fixed.
