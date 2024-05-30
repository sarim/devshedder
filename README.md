# DevShedder

DevShedder is a specialized tool designed to help developers maintain a clean and efficient workspace by automatically detecting and removing unnecessary directories such as `node_modules`, `vendor`, and various log directories from their projects.

## Goal

The primary goal of DevShedder is to free up disk space and reduce clutter in development environments by safely removing directories that are often large and regenerated easily, such as dependencies from Node.js or PHP Composer and log files from Symfony projects.

## Features

- **Selective Cleaning**: Targets specific directories (`node_modules`, `vendor`, `var/log`) based on the project type.
- **Project Detection**: Automatically detects the type of project by looking for key files (`package.json`, `composer.json`, `symfony.lock`).
- **Safe Deletion**: Offers a dry run mode to preview deletions without performing them.
- **User Confirmation**: Requires user confirmation before deletion if not in dry run mode, to prevent accidental data loss.
- **Error Handling**: Gracefully handles errors like permission denials and continues operation, logging issues to the stderr in red color for visibility.

## Usage

### Installation

DevShedder can be installed either by downloading a pre-compiled binary from [github releases](https://github.com/sarim/devshedder/releases) or by building it from the source.

Clone the repository and build the binary using Go:

```bash
git clone https://github.com/sarim/devshedder.git
cd devshedder
go build
```

### Running DevShedder

To clean up your project directory:

```bash
./devshedder /path/to/your/project/directory
```

To perform a dry run and see what would be deleted without actual deletion:

```bash
./devshedder -dry-run /path/to/your/project/directory
```

### Confirmation

If you run without the dry-run option, DevShedder will ask for confirmation before proceeding with deletions:

```bash
Are you sure you want to proceed with deletions? (y/n):
```

Type `y` to proceed or `n` to cancel.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests with your enhancements, or open issues for bugs and feature requests.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Author

DevShedder is developed and maintained by [Sarim Khan](https://github.com/sarim).
