# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- More and better documentation
- `OrderedDict.MustGet()`

### Changed
- `FrozenSet` implementation was modified, avoiding confusion with `Set`.
- Replace build CI job with tests and coverage
- `Dict` has been reimplemented using a slice, instead of a map, because in Go
  not all types can be map's keys (e.g. slices).

### Removed
- Unused method `List.Extend`

## [0.0.1-alpha.1] - 2020-05-23
### Fixed
- Modify GitHub Action steps `Build` and `Test` including all sub-packages.

## [0.0.1-alpha.0] - 2020-05-23
### Added
- Initial implementation of `types` package
- Initial implementation of `pickle` package
- Initial implementation of `pytorch` package
