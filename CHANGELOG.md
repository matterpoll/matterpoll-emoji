# Changelog
All notable changes to this project will be documented in this file.

## 0.1.1 – 2018-01-06
### Fixed
- Set content type to `application/json`([#89](https://github.com/matterpoll/matterpoll-emoji/pull/89))
- Updated dependencies

## 0.1.0 – 2017-08-28
### Breaking changes
- You need add two config parameter to `config.json`:
  - `listen` defines the ip address and port Matterpoll will listen for requests. Default is `localhost:8505`
  - `token` is the mattermost token. You should copy it from Integration > Slash Commands at your Mattermost server.
- Removed `-p` option. Choosing the port is now done via `config.json`

### Added
- Polls are posted as a bot user
- You can use single and double quotation instead of backticks
- We now verify the Mattermost token to make sure requests are valid
- The ip address Matterpoll should listen for request can now we configured via `config.json`
- You can choose the config file via `-c` option
- Go 1.7 to 1.9 are now supported

### Changed
- A lot of unittests were added
- Moved from Mattermost APIv3 to APIv4

### Fixed
- Removed dependency of glide
- Mattermost warning about empty message

## 0.0.2 – 2017-04-19
### Added
- Ci builds via travis

### Changed
- Use 8505 as default port

### Fixed
- Fixed setup documentation

## 0.0.1 – 2017-04-02
- First Version
