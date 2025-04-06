// shadow.fs
#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

out vec4 finalColor;

void main() {
    vec4 texColor = texture(texture0, fragTexCoord);
    
    // If the current pixel is transparent, calculate soft shadow around it
    if (texColor.a == 0.0) {
        float shadow = 0.0;

        // Sample neighboring pixels to see if there's alpha nearby
        float offset = 1.0 / 512.0; // adjust depending on texture size
        for (int x = -2; x <= 2; x++) {
            for (int y = -2; y <= 2; y++) {
                vec2 offsetCoord = fragTexCoord + vec2(x, y) * offset;
                float a = texture(texture0, offsetCoord).a;
                shadow += a;
            }
        }

        shadow = clamp(shadow / 25.0, 0.0, 1.0); // normalize and soften

        finalColor = vec4(0.0, 0.0, 0.0, shadow * 0.4); // black shadow
    } else {
        // Normal texture color for visible pixels
        finalColor = texColor * colDiffuse;
    }
}
