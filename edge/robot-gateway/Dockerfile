FROM ros:humble-ros-base

# Install Python NATS client
RUN pip3 install --no-cache-dir nats-py>=2.6.0

# ROS2 message packages needed by the bridge
RUN apt-get update && apt-get install -y --no-install-recommends \
    ros-humble-geometry-msgs \
    ros-humble-nav-msgs \
    ros-humble-sensor-msgs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY nats_ros2_bridge.py .

# Source ROS2 setup in entrypoint
ENTRYPOINT ["/bin/bash", "-c", "source /opt/ros/humble/setup.bash && python3 /app/nats_ros2_bridge.py"]
