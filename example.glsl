#version 330 core

// Ouput data
out vec4 color;
// Input data
in vec4 gl_FragCoord;

// Uniforms
uniform vec2 iResolution;
uniform float iGlobalTime;

// Texture
uniform sampler2D iChannel0;
uniform sampler2D iChannel1;
uniform sampler2D iChannel2;

void main()
{
  // Output color = red
  //color = vec4(gl_FragCoord.xy / iResolution.xy, 0.5+0.5*sin(iGlobalTime), 1.0);
  vec2 uv = gl_FragCoord.xy / iResolution.xy;
  vec4 s = 3.0 * texture2D(iChannel0, (uv - 0.5) * sin(iGlobalTime) + 0.5) / 4.0;
  vec4 t = vec4(0.0, 1.0, 0.0, 1.0) / 4.0;
  color = s + t;
}
