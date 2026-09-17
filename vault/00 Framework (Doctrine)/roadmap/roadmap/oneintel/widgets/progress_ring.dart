import 'dart:math';
import 'package:flutter/material.dart';

/// Circular progress ring widget.
///
/// Features:
/// - Smooth arc animation
/// - Glow effect when complete
/// - Customizable size, stroke, and colors
class ProgressRing extends StatelessWidget {
  final double progress; // 0.0 to 1.0
  final double size;
  final double strokeWidth;
  final Color color;
  final Color? backgroundColor;
  final Widget? child;

  const ProgressRing({
    super.key,
    required this.progress,
    this.size = 60,
    this.strokeWidth = 6,
    this.color = const Color(0xFFFF0000),
    this.backgroundColor,
    this.child,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: size,
      height: size,
      child: Stack(
        alignment: Alignment.center,
        children: [
          CustomPaint(
            size: Size(size, size),
            painter: _ProgressRingPainter(
              progress: progress.clamp(0.0, 1.0),
              strokeWidth: strokeWidth,
              color: color,
              backgroundColor: backgroundColor ?? Colors.white10,
            ),
          ),
          if (child != null) child!,
        ],
      ),
    );
  }
}

class _ProgressRingPainter extends CustomPainter {
  final double progress;
  final double strokeWidth;
  final Color color;
  final Color backgroundColor;

  _ProgressRingPainter({
    required this.progress,
    required this.strokeWidth,
    required this.color,
    required this.backgroundColor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = (size.width - strokeWidth) / 2;

    // Background ring
    final backgroundPaint = Paint()
      ..color = backgroundColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;

    canvas.drawCircle(center, radius, backgroundPaint);

    // Progress arc
    if (progress > 0) {
      final progressPaint = Paint()
        ..color = color
        ..style = PaintingStyle.stroke
        ..strokeWidth = strokeWidth
        ..strokeCap = StrokeCap.round;

      final sweepAngle = 2 * pi * progress;

      canvas.drawArc(
        Rect.fromCircle(center: center, radius: radius),
        -pi / 2, // Start from top
        sweepAngle,
        false,
        progressPaint,
      );

      // Glow effect for completed
      if (progress >= 1.0) {
        final glowPaint = Paint()
          ..color = color.withValues(alpha: 0.3)
          ..style = PaintingStyle.stroke
          ..strokeWidth = strokeWidth * 2
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 6);

        canvas.drawCircle(center, radius, glowPaint);
      }
    }
  }

  @override
  bool shouldRepaint(covariant _ProgressRingPainter oldDelegate) {
    return oldDelegate.progress != progress ||
        oldDelegate.color != color ||
        oldDelegate.strokeWidth != strokeWidth;
  }
}

/// Large progress ring with percentage label.
class LargeProgressRing extends StatelessWidget {
  final double progress;
  final String? label;
  final String? subtitle;
  final Color color;

  const LargeProgressRing({
    super.key,
    required this.progress,
    this.label,
    this.subtitle,
    this.color = const Color(0xFFFF0000),
  });

  @override
  Widget build(BuildContext context) {
    final percentage = (progress * 100).toInt();

    return ProgressRing(
      progress: progress,
      size: 120,
      strokeWidth: 10,
      color: color,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            '$percentage%',
            style: const TextStyle(
              color: Colors.white,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          if (label != null)
            Text(
              label!,
              style: const TextStyle(
                color: Colors.white54,
                fontSize: 12,
              ),
            ),
          if (subtitle != null)
            Text(
              subtitle!,
              style: TextStyle(
                color: color,
                fontSize: 10,
              ),
            ),
        ],
      ),
    );
  }
}

/// Small progress ring for list items.
class MiniProgressRing extends StatelessWidget {
  final double progress;
  final Color? color;

  const MiniProgressRing({
    super.key,
    required this.progress,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    final ringColor = color ?? _getProgressColor(progress);

    return ProgressRing(
      progress: progress,
      size: 32,
      strokeWidth: 3,
      color: ringColor,
      child: Text(
        '${(progress * 100).toInt()}',
        style: TextStyle(
          color: ringColor,
          fontSize: 9,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Color _getProgressColor(double progress) {
    if (progress >= 0.8) return Colors.green;
    if (progress >= 0.6) return Colors.yellow.shade700;
    if (progress >= 0.4) return Colors.orange;
    return Colors.red;
  }
}
