# Create distribution directory
$distDir = "dist/native-radio"
New-Item -ItemType Directory -Force -Path $distDir

# Build the executable
go build -o "$distDir/native-radio.exe"

# Copy required assets
Copy-Item -Path "fonts" -Destination "$distDir/fonts" -Recurse
Copy-Item -Path "shaders" -Destination "$distDir/shaders" -Recurse
Copy-Item -Path "sprites" -Destination "$distDir/sprites" -Recurse
Copy-Item -Path "masks" -Destination "$distDir/masks" -Recurse
Copy-Item -Path "resources" -Destination "$distDir/resources" -Recurse

# Create ZIP archive
# Compress-Archive -Path "$distDir/*" -DestinationPath "native-radio-windows.zip" -Force