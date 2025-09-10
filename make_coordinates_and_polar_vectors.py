import numpy as np
import matplotlib.pyplot as plt

# Parameters
N = 100  # number of coordinates
np.random.seed(42)

# Step 1: Generate random polar vectors (r, theta)
r = np.random.randint(1, 6, size=(N-1))  # random step lengths 1–5
theta = np.random.uniform(0, 2*np.pi, size=(N-1))  # random angles in radians

# Step 2: Convert polar -> Cartesian
dx = r * np.cos(theta)
dy = r * np.sin(theta)
vectors = np.vstack([dx, dy]).T

# Step 3: Build coordinates from vectors
start = np.array([0, 0])
coordinates = np.vstack([start, start + np.cumsum(vectors, axis=0)])

# Print both forms
print("Vectors (Polar): r, θ (radians):")
for ri, ti in zip(r, theta):
    print(f"r={ri:.2f}, θ={ti:.2f} rad")

print("\nCoordinates:")
for c in coordinates:
    print(f"({c[0]:.2f}, {c[1]:.2f})")

# Step 4: Plot path with vectors
plt.figure(figsize=(6,6))
plt.plot(coordinates[:,0], coordinates[:,1], 'o-', label="Path")
for i in range(len(vectors)):
    plt.arrow(coordinates[i,0], coordinates[i,1], vectors[i,0], vectors[i,1],
              head_width=0.2, length_includes_head=True, color='red')
plt.xlabel("X")
plt.ylabel("Y")
plt.title(f"Random Path with {N} Coordinates (Polar Vectors)")
plt.grid(True)
plt.legend()
plt.axis("equal")
plt.show()

"""
Vectors (Polar): r, θ (radians):
r=4.00, θ=0.73 rad
r=5.00, θ=5.42 rad
r=3.00, θ=3.92 rad
r=5.00, θ=2.08 rad
r=5.00, θ=0.40 rad
r=2.00, θ=1.95 rad
r=3.00, θ=2.04 rad
r=3.00, θ=4.58 rad
r=3.00, θ=4.01 rad
r=5.00, θ=5.57 rad
r=4.00, θ=2.97 rad
r=3.00, θ=0.75 rad
r=5.00, θ=4.48 rad
r=2.00, θ=4.78 rad
r=4.00, θ=3.53 rad
r=2.00, θ=4.84 rad
r=4.00, θ=3.10 rad
r=5.00, θ=3.28 rad
r=1.00, θ=2.69 rad
r=4.00, θ=0.16 rad
r=2.00, θ=0.68 rad
r=5.00, θ=0.20 rad
r=4.00, θ=4.00 rad
r=1.00, θ=1.98 rad
r=1.00, θ=3.20 rad
r=3.00, θ=5.70 rad
r=3.00, θ=1.57 rad
r=2.00, θ=2.58 rad
r=4.00, θ=4.75 rad
r=4.00, θ=1.44 rad
r=3.00, θ=0.48 rad
r=4.00, θ=1.82 rad
r=4.00, θ=1.01 rad
r=1.00, θ=5.84 rad
r=3.00, θ=5.08 rad
r=5.00, θ=3.98 rad
r=3.00, θ=5.48 rad
r=5.00, θ=5.05 rad
r=1.00, θ=1.17 rad
r=2.00, θ=5.61 rad
r=4.00, θ=3.39 rad
r=1.00, θ=5.07 rad
r=4.00, θ=5.63 rad
r=2.00, θ=2.00 rad
r=2.00, θ=0.69 rad
r=1.00, θ=1.43 rad
r=2.00, θ=2.68 rad
r=5.00, θ=5.14 rad
r=2.00, θ=5.41 rad
r=4.00, θ=0.04 rad
r=4.00, θ=3.21 rad
r=4.00, θ=2.62 rad
r=4.00, θ=1.40 rad
r=5.00, θ=0.75 rad
r=3.00, θ=2.12 rad
r=1.00, θ=5.92 rad
r=4.00, θ=2.03 rad
r=2.00, θ=3.26 rad
r=4.00, θ=4.42 rad
r=2.00, θ=2.28 rad
r=2.00, θ=6.11 rad
r=4.00, θ=6.05 rad
r=5.00, θ=1.58 rad
r=2.00, θ=3.12 rad
r=2.00, θ=1.89 rad
r=4.00, θ=1.79 rad
r=2.00, θ=0.23 rad
r=2.00, θ=3.83 rad
r=4.00, θ=3.16 rad
r=4.00, θ=0.32 rad
r=1.00, θ=1.75 rad
r=5.00, θ=5.71 rad
r=5.00, θ=1.51 rad
r=2.00, θ=0.91 rad
r=5.00, θ=3.08 rad
r=2.00, θ=6.19 rad
r=1.00, θ=1.52 rad
r=4.00, θ=4.22 rad
r=4.00, θ=4.79 rad
r=4.00, θ=1.49 rad
r=5.00, θ=4.58 rad
r=1.00, θ=2.31 rad
r=5.00, θ=3.97 rad
r=5.00, θ=3.98 rad
r=1.00, θ=3.37 rad
r=1.00, θ=0.57 rad
r=1.00, θ=5.25 rad
r=1.00, θ=2.02 rad
r=4.00, θ=1.17 rad
r=3.00, θ=0.26 rad
r=3.00, θ=3.71 rad
r=1.00, θ=4.26 rad
r=3.00, θ=0.10 rad
r=3.00, θ=3.22 rad
r=1.00, θ=1.42 rad
r=3.00, θ=4.05 rad
r=5.00, θ=1.10 rad
r=2.00, θ=4.34 rad
r=2.00, θ=2.43 rad

Coordinates:
(0.00, 0.00)
(2.99, 2.66)
(6.25, -1.13)
(4.10, -3.23)
(1.67, 1.14)
(6.28, 3.09)
(5.53, 4.94)
(4.16, 7.61)
(3.78, 4.64)
(1.83, 2.35)
(5.63, -0.90)
(1.69, -0.20)
(3.88, 1.84)
(2.74, -3.02)
(2.87, -5.02)
(-0.83, -6.52)
(-0.57, -8.50)
(-4.57, -8.35)
(-9.52, -9.06)
(-10.42, -8.62)
(-6.47, -7.98)
(-4.91, -6.73)
(-0.01, -5.75)
(-2.62, -8.77)
(-3.02, -7.85)
(-4.02, -7.91)
(-1.51, -9.55)
(-1.49, -6.55)
(-3.19, -5.49)
(-3.05, -9.48)
(-2.52, -5.52)
(0.14, -4.12)
(-0.85, -0.25)
(1.27, 3.15)
(2.17, 2.72)
(3.24, -0.08)
(-0.10, -3.80)
(1.97, -5.97)
(3.63, -10.69)
(4.02, -9.77)
(5.58, -11.02)
(1.70, -11.99)
(2.05, -12.93)
(5.23, -15.36)
(4.40, -13.54)
(5.94, -12.26)
(6.08, -11.27)
(4.29, -10.39)
(6.36, -14.94)
(7.64, -16.48)
(11.64, -16.30)
(7.65, -16.57)
(4.17, -14.59)
(4.87, -10.65)
(8.52, -7.23)
(6.95, -4.67)
(7.88, -5.02)
(6.11, -1.44)
(4.12, -1.67)
(2.96, -5.50)
(1.65, -3.99)
(3.62, -4.34)
(7.51, -5.28)
(7.45, -0.28)
(5.45, -0.24)
(4.82, 1.66)
(3.95, 5.56)
(5.90, 6.02)
(4.36, 4.75)
(0.36, 4.68)
(4.15, 5.95)
(3.97, 6.94)
(8.16, 4.21)
(8.49, 9.20)
(9.72, 10.78)
(4.73, 11.11)
(6.72, 10.93)
(6.77, 11.93)
(4.89, 8.40)
(5.18, 4.41)
(5.49, 8.40)
(4.81, 3.45)
(4.14, 4.18)
(0.77, 0.49)
(-2.57, -3.23)
(-3.55, -3.45)
(-2.71, -2.92)
(-2.20, -3.78)
(-2.63, -2.87)
(-1.07, 0.81)
(1.83, 1.57)
(-0.69, -0.05)
(-1.13, -0.95)
(1.85, -0.63)
(-1.14, -0.86)
(-0.99, 0.13)
(-2.83, -2.25)
(-0.54, 2.20)
(-1.27, 0.34)
(-2.78, 1.64)
"""