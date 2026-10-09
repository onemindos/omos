# OpenChoreo gateway compatibility

The pinned 1.3.1 chart emits an assistant HTTPRoute rule name unsupported by the
installed Gateway API schema. Disable that chart route and supply the same paths,
timeouts and backend through this compatible override. The API route remains
chart-owned. Remove this override only after installed schema support and a
successful server-side validation of the rendered chart.

The chart's HTTPS listener requires certificateRefs, which were unset. The
internal gateway uses HTTP and the existing onemind-private GatewayClass keeps its
Service ClusterIP. This adds no public exposure and preserves upstream
authentication. Certificate/mesh integration remains an acceptance check before
exposing this gateway beyond the cluster.
