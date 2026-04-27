# Assets tests

Scope: asset kind metadata, resource cache concurrency/lifecycle, manager wrappers and main-thread queue, asset generator config/scanning/naming/template behavior.

### Kind metadata (`assets/asset_kind_test.go`, `assets/manager_test.go`)

- `TestKindDefaultDirKnownAndCustom`: verifies default directory mapping for built-in kinds and custom plural fallback.
- `TestKindPluralRules`: verifies pluralization rules for built-ins and common word endings.
- `TestKindIsKnown`: verifies built-in vs custom kind detection.
- `TestKindDefaultExtensions`: verifies built-in extension sets and nil extension list for custom kinds.
- `TestDetectKind`: verifies extension-to-kind mapping, including unknown fallback.
- `TestManagerDetectKindUsesBuiltInExtensions`: verifies case-insensitive extension handling.

### Cache concurrency (`assets/manager_test.go`)

- `TestResourceCacheLoadConcurrentSingleLoader`: verifies concurrent same-key load deduplicates loader call and returns shared handle pointer.

### Cache lifecycle and reload semantics (`assets/manager_test.go`)

- `TestResourceCacheReloadErrorKeepsOldData`: failed reload must preserve previous data.
- `TestResourceHandleUnloadAndReloadReattach`: unload detaches+zeros handle; reload reattaches same handle pointer.
- `TestResourceSafeRead`: `SafeRead` callback observes current value.
- `TestResourceCacheClearUnloadsAllAndZeros`: `Clear` unloads all, clears keys, zeroes handles.
- `TestResourceHandleReloadFailureWhenDetachedKeepsCacheEmpty`: failed detached reload must not recreate cache entry.

### Manager API and threading contracts (`assets/manager_test.go`)

- `TestFontLoaderRejectsInvalidSizeKey`: validates font size key parsing errors.
- `TestManagerReloadAllCallsErrorCallback`: validates reload-all error callback flow and nil callback tolerance.
- `TestManagerWrapperMethodsAndBulkOps`: validates typed `Get*`/`Reload*`/`Unload*`, accessors, key listing, clear-all, close.
- `TestManagerRunOnMainQueuePath`: validates queued execution path for non-owner goroutines.
- `TestManagerCloseIsIdempotentAndStopsQueueing`: validates idempotent close and post-close queue no-op.
- `TestManagerKeysSwitchAllKinds`: validates `Keys(kind)` switch coverage and unknown-kind behavior.
- `TestNewManagerLoaderClosuresErrorPaths`: validates initialization and loader error paths for missing assets/shader key format.

### Generator happy paths (`assets/generate_test.go`)

- `TestGenerateBasic`: verifies default generation emits package/type/constants for discovered textures.
- `TestGenerateWithExtensions`: verifies explicit kind extension mapping and exclusion of non-matching files.
- `TestFlatMode`: verifies flat mode emits one shared type.
- `TestDryRunWritesToWriter`: verifies dry-run writes generated output to provided writer.
- `TestGenerateDryRunAndDump`: verifies `GenerateDryRun` and `Config.Dump` helper behavior.

### Generator config merge and defaults (`assets/generate_test.go`)

- `TestLoadFileConfigReturnsProperConfig`: validates config JSON field mapping.
- `TestLoadFileConfigWithKinds`: validates kinds array parsing.
- `TestMergeWithFile`: validates merge precedence rules for strings/bools/slices/kinds.
- `TestGenerateMergesConfigWithFileConfig`: validates programmatic config overrides file config.
- `TestGenerateUsesFileConfigWhenProgrammaticIsZero`: validates file config fills missing programmatic fields.
- `TestApplyStringDefaults`: validates string defaults and non-overwrite behavior.
- `TestReadFileConfigBranches`: validates file config read error branches.
- `TestNormalizeKindsBranches`: validates kind normalization and single-root behavior.

### Generator scanning and filtering (`assets/generate_test.go`)

- `TestMatchInclExcl`: validates include/exclude precedence and matching semantics.
- `TestScanDir`: validates recursive behavior, glob filtering, hidden file skip, collision branch.
- `TestScanDirWithFilters`: validates include/exclude combinations.
- `TestScanDirStripPrefix`: validates prefix stripping in emitted paths.
- `TestTransformPath`: validates absolute/relative normalization and optional prefix strip.
- `TestSkipForHidden`: validates hidden-file filtering.
- `TestSkipForGlob`: validates glob filtering.

### Generator naming, template, and error paths (`assets/generate_test.go`)

- `TestNameCollisionIncludesBothPaths`: collision errors include both source paths.
- `TestGenerateRejectsInvalidNaming`: invalid naming style returns validation error.
- `TestGenerateConfigFileParseError`: malformed config file returns parse error.
- `TestGenerateTemplateFileReadError`: missing template file returns read error.
- `TestGenerateTemplateParseError`: invalid template syntax returns parse error.
- `TestNamerAllStylesAndPrefix`: validates naming styles, sanitization, prefix behavior.
- `TestHelperFunctions`: validates utility helpers (`extensionsToGlob`, `capitalize`, `isVowel`).
- `TestWarnUnknownKindsVerbose`: validates warning output for custom kinds in verbose mode.
- `TestScanDirCollision`: validates explicit scanner collision error path.