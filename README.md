# VersionUp

The VersionUp project was created to simplify the management of project versions (tags) without the need for manual intervention.

## Environment Variables

| Variables         | Description                                                                                           |
| ----------------  | ---------------------------------------------------------------------------------------------------   |
| `ENVIROMENT`      | Environment to update the version (dev, hmg, prd)                                                     |
| `TYPE_VERSION`    | Used to represent the meaning of the realized version (major, minor, patch)                           |
| `DESCRIPTION_TAG` | Set a name for the prefix of a tag                                                                    |
| `REPOSITORY_URL`  | Repository URL                                                                                        |
| `USER_NAME`       | Git username                                                                                          |
| `USER_EMAIL`      | Git user email                                                                                        |
| `SSH_KNOWN_HOSTS` | Used to specify the path of the "known_hosts" file containing information about trusted remote hosts  |
| `SSH_AUTH_SOCK`   | Used to set the path of the SSH authentication agent socket                                           |

`TYPE_VERSION`:
- `major`: Involves major restructurings, new features, or API design changes
- `minor`: Incremental updates that do not break compatibility with previous usage
- `patch`: Updates that fix specific issues without adding new features or changing APIs

## Configuration of your machine to define volumes
Access the `~/.ssh` folder and create a subfolder inside it called `/ssh_versionup`. 
Once done, copy the `id_rsa` and `known_hosts` files into the `/ssh_versionup` folder. 
Run the following command from the `~/.ssh` folder:

```bash
ssh-keygen -p -f /ssh_versionup/id_rsa
```

## Volumes

| From                                 | Mount                      | Description                                                                                                              |
| -------------------------            | -------------------------  |---------------------------------------------------------------------------------------------------   |
| `$(pwd)/versionup.yml`               | /app/versionup_old.yml     | It's important to include the versionup.yml file in your project as it provides the necessary version information for the tool                                  |
| `~/.ssh/ssh_versionup/id_rsa`        | /ssh/id_rsa                | The id_rsa file is a private key file used by the SSH (Secure Shell) protocol, which is located in the ~/.ssh/ssh_versionup folder                           |
| `~/.ssh/ssh_versionup/known_hosts`   | /ssh/known_hosts           | The known_hosts file is a text file used by the SSH (Secure Shell) protocol to store information about remote hosts that you have previously connected to via SSH |
| `$SSH_AUTH_SOCK`                     | /ssh-agent.sock            | The $SSH_AUTH_SOCK environment variable is used by the operating system to define the location of the SSH authentication agent socket                           |

Note: The paths declared in `From` were made from the Linux OS, so depending on the OS they may suffer a little change

## Configuration File
To use the tool, create a `versionup.yml` file with the following structure:

```yaml
version:
  dev: 1.0.0
  hmg: 1.0.0
  prd: 1.0.0
```

- `dev`: Represents the version in the development environment
- `hmg`: Represents the version in the homologation environment
- `prd`: Represents the version in the production environment

Replace the default value "1.0.0" with the appropriate version number for each environment in your project.

This command will prompt you for the passphrase of this SSH private key and then ask for the new passphrase, which should be left blank.

## Usage
To run the tool via Docker, run the following command:
```bash
docker run -e ENVIROMENT=dev -e TYPE_VERSION=patch -e DESCRIPTION_TAG=dev_v -e REPOSITORY_URL=https://github.com/jose/example -e USER_NAME=Jose -e USER_EMAIL=JoseDasCoves@gmail.com -e SSH_KNOWN_HOSTS=/ssh/known_hosts -e SSH_AUTH_SOCK=/ssh-agent.sock -v "$(pwd)/versionup.yml:/app/versionup_old.yml" -v ~/.ssh/ssh_versionup/id_rsa:/ssh/id_rsa -v ~/.ssh/ssh_versionup/known_hosts:/ssh/known_hosts -v "$SSH_AUTH_SOCK:/ssh-agent.sock" versionup
```

After running the above command, don't forget to update your local project, with the command:
```bash
git pull origin
```