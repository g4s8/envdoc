# Environment Variables

## Config

Config is an example configuration structure.
It is used to generate documentation for the configuration
using the commands below.

| env variable | type | description | options |
|--------------|------|-------------|---------|
| HOST | `[]string` | Hosts name of hosts to listen on. | (separated by `;`, **required**) |
| PORT | `int` | Port to listen on. | (**required**, non-empty) |
| DEBUG | `bool` | Debug mode enabled. | (default: `false`) |
| PREFIX | `string` | Prefix for something. |  |

