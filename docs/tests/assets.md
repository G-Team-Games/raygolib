# Assets tests

Scope: asset kind metadata, resource cache concurrency/lifecycle, manager wrappers and main-thread queue, asset generator config/scanning/naming/template behavior.

## `assets/asset_kind_test.go`

- `TestKindDefaultDirKnownAndCustom`: verifies default directory mapping for built-in kinds and custom plural fallback.
- `TestKindPluralRules`: verifies pluralization rules for built-ins and common word endings.
- `TestKindIsKnown`: verifies built-in vs custom kind detection.
- `TestKindDefaultExtensions`: verifies built-in extension sets and nil extension list for custom kinds.

## `assets/manager_test.go`

- `TestResourceCacheLoadConcurrentSingleLoader`: verifies concurrent same-key load deduplicates loader call and returns shared handle pointer.
- `TestResourceCacheReloadErrorKeepsOldData`: verifies failed reload preserves previously loaded data.
- `TestResourceHandleUnloadAndReloadReattach`: verifies unload detaches+zeros handle and later reload reattaches same handle.
- `TestResourceSafeRead`: verifies `SafeRead` callback reads consistent resource value.
- `TestResourceCacheClearUnloadsAllAndZeros`: verifies `Clear` unloads every resource, empties keys, zeroes handles.
- `TestResourceHandleReloadFailureWhenDetachedKeepsCacheEmpty`: verifies failed reload of detached handle does not insert cache entry.
- `TestDetectKind`: verifies extension-to-kind mapping, including unknown fallback.
- `TestManagerDetectKindUsesBuiltInExtensions`: verifies case-insensitive built-in extension mapping.
- `TestFontLoaderRejectsInvalidSizeKey`: verifies font loader rejects invalid/non-positive size suffix.
- `TestManagerReloadAllCallsErrorCallback`: verifies `ReloadAll` aggregates reload errors via callback and tolerates nil callback.
- `TestManagerWrapperMethodsAndBulkOps`: verifies typed getters/accessors, reload wrappers, key listing, unload wrappers, clear-all, close.
- `TestManagerRunOnMainQueuePath`: verifies non-owner goroutine path queues work and `Tick` executes queued callback.
- `TestManagerCloseIsIdempotentAndStopsQueueing`: verifies close idempotency and post-close queue no-op.
- `TestManagerKeysSwitchAllKinds`: verifies `Keys(kind)` switch covers all built-in kinds and returns nil for unknown.
- `TestNewManagerLoaderClosuresErrorPaths`: verifies `NewManager` initializes all caches and loaders return expected errors on missing assets/bad shader key.

## `assets/generate_test.go`

- `TestGenerateBasic`: verifies basic generation emits package/type/constants for discovered textures.
- `TestGenerateWithExtensions`: verifies explicit kind extension mapping and exclusion of non-matching files.
- `TestFlatMode`: verifies flat mode emits one shared type with configured type name.
- `TestDryRunWritesToWriter`: verifies dry-run output path writes generated source to writer.
- `TestNameCollisionIncludesBothPaths`: verifies collision error includes both conflicting source paths.
- `TestGenerateDryRunAndDump`: verifies helper wrappers (`GenerateDryRun`, `Config.Dump`) output expected content.
- `TestGenerateRejectsInvalidNaming`: verifies invalid naming style returns validation error.
- `TestGenerateConfigFileParseError`: verifies malformed config file returns parse error.
- `TestGenerateTemplateFileReadError`: verifies missing template file returns read error.
- `TestGenerateTemplateParseError`: verifies invalid template syntax returns parse error.
- `TestHelperFunctions`: verifies helper utilities (`extensionsToGlob`, `capitalize`, `isVowel`).
- `TestNamerAllStylesAndPrefix`: verifies naming styles, sanitization, prefix handling.
- `TestReadFileConfigBranches`: verifies `loadFileConfig` error branches (empty path, missing path).
- `TestNormalizeKindsBranches`: verifies kind normalization defaults and single-root behavior.
- `TestWarnUnknownKindsVerbose`: verifies unknown kind warning output to stderr.
- `TestLoadFileConfigReturnsProperConfig`: verifies JSON config mapping to `Config` fields.
- `TestLoadFileConfigWithKinds`: verifies config file kind entries parse into expected kind definitions.
- `TestMergeWithFile`: verifies merge precedence for strings, booleans, lists, and kinds.
- `TestGenerateMergesConfigWithFileConfig`: verifies programmatic config overrides file config when both set.
- `TestGenerateUsesFileConfigWhenProgrammaticIsZero`: verifies file config fills zero-valued programmatic config.
- `TestApplyStringDefaults`: verifies default strings and no-op on pre-set values.
- `TestMatchInclExcl`: verifies include/exclude matching semantics and precedence.
- `TestScanDir`: verifies recursive and non-recursive scanning, glob filtering, hidden file skip, collision behavior.
- `TestScanDirWithFilters`: verifies include/exclude filter combinations.
- `TestScanDirStripPrefix`: verifies path prefix stripping in emitted asset paths.
- `TestScanDirCollision`: verifies explicit collision error path in scanner.
- `TestTransformPath`: verifies absolute/relative path normalization and optional prefix strip.
- `TestSkipForHidden`: verifies hidden file detection.
- `TestSkipForGlob`: verifies glob skip behavior.