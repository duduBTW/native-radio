#version 330

const int samples = 55;
const int LOD = 3;
const int sLOD = 1 << LOD;

uniform sampler2D texture0;
uniform vec2 resolution;
uniform vec2 textureResolution;
uniform vec2 mouse;
uniform float blurRadius;

out vec4 fragColor;

float gaussian(vec2 i, float sigma) {
    return exp(-0.5 * dot(i /= sigma, i)) / (6.28 * sigma * sigma);
}

vec4 blur(sampler2D sp, vec2 U, vec2 scale, float sigma) {
    vec4 O = vec4(0.0);  
    int s = samples / sLOD;

    for (int i = 0; i < s * s; i++) {
        vec2 d = vec2(i % s, i / s) * float(sLOD) - float(samples) / 2.0;
        O += gaussian(d, sigma) * textureLod(sp, U + scale * d, float(LOD));
    }

    return O / O.a;
}

void main() {
    vec2 uv = gl_FragCoord.xy / resolution;
    uv.y = 1.0 - uv.y;

    float dist = distance(gl_FragCoord.xy, mouse);
    float factor = smoothstep(0.0, blurRadius, dist);
    float sigma = mix(float(samples)*0.01, float(samples)*0.2, factor);

    fragColor = blur(texture0, uv, 1.0 / textureResolution, sigma);
}