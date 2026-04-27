# Collision 3D tests

Scope: contact detection, penetration/normal contracts, distance invariants, raycast correctness, resolver behavior, collider draw helpers.

## `physics/collision/3d/collision_test.go`

- `TestBoxVsBoxContact`: verifies overlap hit with positive penetration and separated miss for box-box.
- `TestBoxHelpers`: verifies box center/min/max helper math.
- `TestRaycastBox`: verifies box raycast hit + hit point population.
- `TestCylinderVsBoxContactAndResolve`: verifies cylinder-box contact and `ResolveMTV` moves collider.
- `TestPlaneAxisNormal`: verifies axis enum returns expected plane normal.
- `TestCylinderCenter`: verifies cylinder center math.
- `TestPlaneCenter`: verifies center calculations for X/Y/Z planes.
- `TestContract_CollideNormalPointsFromBToA`: verifies normal direction contract (from B toward A).
- `TestContract_TouchingIsHitWithZeroPenetration`: verifies touching counts as hit with zero penetration.
- `TestDistance_BoxVsBox_OverlapIsZero`: verifies distance is zero for overlapping boxes.
- `TestDistance_CylinderVsCylinder_OverlapIsZero`: verifies distance is zero for overlapping cylinders.
- `TestDistance_PointVsCylinder_OutsideCircleNonZero`: verifies positive distance for outside point.
- `TestDistance_SymmetryForSupportedPairs`: verifies distance symmetry across supported pair families.
- `TestDistance_UnsupportedPairReturnsInf`: verifies unsupported pair distance returns `+Inf`.
- `TestResolveMTV_ReducesOverlap_BoxBoxInOneStep`: verifies one-step MTV resolve reduces overlap to touching state.
- `TestCollide_ReverseOrderNormalsAreOpposite_WhenBothSupported`: verifies reverse-order collisions keep penetration and invert normal.
- `TestCylinderVsHorizontalPlane_MinimalTranslationDirection`: verifies nearest MTV direction and expected penetration for cylinder-plane.
- `TestCollide_OrderIndependence_AllSupportedUnorderedPairs`: verifies order invariants across matrix of supported pairs.
- `TestRaycastCylinder_SideHit`: verifies side-hit intersection point/normal.
- `TestRaycastCylinder_TopHit`: verifies top-cap intersection point/normal.
- `TestRaycastCylinder_BottomHit`: verifies bottom-cap intersection point/normal.
- `TestRaycastCylinder_Miss`: verifies miss case returns no hit.
- `TestRaycastCylinder_TallCylinderTopHit`: verifies tall-cylinder top-cap raycast.
- `TestRaycastPlane_NormalMatchesAxis`: verifies raycast hit normals follow axis orientation variants.
- `TestCollisionMatrix_BoundaryAndReverseOrder`: verifies broad matrix of touch/overlap/separation scenarios + reverse-order invariants.
- `TestDistanceInvariants_Matrix`: verifies finite/non-negative/symmetric distance and zero-vs-positive behavior across matrix.

### Subtests with broad scenario matrices

- `TestCollisionMatrix_BoundaryAndReverseOrder/*`: per-scenario assertions for hit expectation, penetration sign, reverse-order normal inversion.
- `TestDistanceInvariants_Matrix/*`: per-scenario assertions for finite distance, symmetry, overlap zero, separation positive.

## `physics/collision/3d/resolver_test.go`

- `TestResolveMultiMTV`: verifies multi-collider iterative resolve converges when sufficient iteration budget exists.

## `physics/collision/3d/draw_colliders_test.go`

- `TestPlaneDrawBox_AppliesThicknessOnFlatAxis`: verifies debug draw helper injects default thickness on plane flat axis.
- `TestPlaneDrawBox_CenterUsesExpandedSize`: verifies draw center computed from expanded size.

## Why this category matters

- Guards engine physics contracts where subtle normal/penetration bugs cause visible gameplay defects.
- Guards order-independence and symmetry invariants needed for deterministic behavior.
- Guards raycast and debug draw correctness used by gameplay and tooling.
