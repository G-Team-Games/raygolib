# Collision 3D tests

Scope: contact detection, penetration/normal contracts, distance invariants, raycast correctness, resolver behavior, collider draw helpers.

### Geometry helpers (`physics/collision/3d/collision_test.go`)

- `TestBoxHelpers`: verifies box center/min/max helper math.
- `TestPlaneAxisNormal`: verifies axis enum returns expected plane normal.
- `TestCylinderCenter`: verifies cylinder center math.
- `TestPlaneCenter`: verifies center calculations for X/Y/Z planes.

### Contact contracts (`physics/collision/3d/collision_test.go`)

- `TestBoxVsBoxContact`: overlap hit with positive penetration and separated miss for box-box.
- `TestCylinderVsBoxContactAndResolve`: cylinder-box hit and movement after `ResolveMTV`.
- `TestContract_CollideNormalPointsFromBToA`: normal direction contract (from B toward A).
- `TestContract_TouchingIsHitWithZeroPenetration`: touching counts as hit with zero penetration.
- `TestCollide_ReverseOrderNormalsAreOpposite_WhenBothSupported`: reverse order keeps penetration, flips normal.
- `TestCylinderVsHorizontalPlane_MinimalTranslationDirection`: nearest MTV direction and penetration for cylinder-plane.
- `TestCollide_OrderIndependence_AllSupportedUnorderedPairs`: order invariants across supported collider pairs.

### Distance invariants (`physics/collision/3d/collision_test.go`)

- `TestDistance_BoxVsBox_OverlapIsZero`: overlapping boxes have zero distance.
- `TestDistance_CylinderVsCylinder_OverlapIsZero`: overlapping cylinders have zero distance.
- `TestDistance_PointVsCylinder_OutsideCircleNonZero`: outside point has positive distance.
- `TestDistance_SymmetryForSupportedPairs`: supported pairs are symmetric.
- `TestDistance_UnsupportedPairReturnsInf`: unsupported pair distance returns `+Inf`.
- `TestDistanceInvariants_Matrix`: matrix coverage for finite/non-negative/symmetric/zero-vs-positive invariants.

### Raycast behavior (`physics/collision/3d/collision_test.go`)

- `TestRaycastBox`: box raycast hit + hit point population.
- `TestRaycastCylinder_SideHit`: side-hit intersection point/normal.
- `TestRaycastCylinder_TopHit`: top-cap intersection point/normal.
- `TestRaycastCylinder_BottomHit`: bottom-cap intersection point/normal.
- `TestRaycastCylinder_Miss`: miss case returns no hit.
- `TestRaycastCylinder_TallCylinderTopHit`: top-cap hit on tall cylinder.
- `TestRaycastPlane_NormalMatchesAxis`: hit normals follow each plane axis orientation.

### Resolver behavior (`physics/collision/3d/collision_test.go`, `physics/collision/3d/resolver_test.go`)

- `TestResolveMTV_ReducesOverlap_BoxBoxInOneStep`: one-step resolve reduces overlap to touching.
- `TestResolveMultiMTV`: iterative multi-collider resolve converges with sufficient iterations.

### Draw helper behavior (`physics/collision/3d/draw_colliders_test.go`)

- `TestPlaneDrawBox_AppliesThicknessOnFlatAxis`: draw helper injects thickness on plane flat axis.
- `TestPlaneDrawBox_CenterUsesExpandedSize`: draw center computed from expanded size.

### Matrix regression suites

- `TestCollisionMatrix_BoundaryAndReverseOrder`: scenario matrix for touch/overlap/separation + reverse-order invariants.
- `TestCollisionMatrix_BoundaryAndReverseOrder/*`: per-scenario hit expectation, penetration sign, opposite normals.
- `TestDistanceInvariants_Matrix/*`: per-scenario finite/non-negative/symmetric/zero-vs-positive checks.
