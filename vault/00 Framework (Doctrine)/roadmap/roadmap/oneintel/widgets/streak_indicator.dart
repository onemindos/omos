import 'package:flutter/material.dart';

/// Animated streak indicator widget.
///
/// Shows different tiers based on streak length:
/// - 1-6 days: Spark (orange)
/// - 7-29 days: Warm (deep orange)
/// - 30-99 days: On Fire (red-orange with glow)
/// - 100+ days: Legendary (gold with glow)
class StreakIndicator extends StatelessWidget {
  final int streak;
  final bool showLabel;
  final double size;

  const StreakIndicator({
    super.key,
    required this.streak,
    this.showLabel = true,
    this.size = 24,
  });

  @override
  Widget build(BuildContext context) {
    if (streak == 0) {
      return const SizedBox.shrink();
    }

    final tier = _getStreakTier(streak);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        _buildFlameIcon(tier),
        if (showLabel) ...[
          const SizedBox(width: 4),
          Text(
            '$streak',
            style: TextStyle(
              color: tier.color,
              fontSize: size * 0.6,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ],
    );
  }

  Widget _buildFlameIcon(StreakTier tier) {
    return Container(
      width: size,
      height: size,
      decoration: tier.hasGlow
          ? BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: tier.color.withValues(alpha: 0.4),
                  blurRadius: size * 0.5,
                  spreadRadius: size * 0.1,
                ),
              ],
            )
          : null,
      child: Icon(
        Icons.local_fire_department,
        color: tier.color,
        size: size,
      ),
    );
  }

  StreakTier _getStreakTier(int streak) {
    if (streak >= 100) {
      return StreakTier.legendary;
    } else if (streak >= 30) {
      return StreakTier.fire;
    } else if (streak >= 7) {
      return StreakTier.warm;
    } else {
      return StreakTier.spark;
    }
  }
}

/// Streak tier levels with associated colors and effects.
enum StreakTier {
  spark(color: Colors.orange, hasGlow: false, label: 'Spark'),
  warm(color: Color(0xFFFF6B35), hasGlow: false, label: 'Warm'),
  fire(color: Color(0xFFFF4500), hasGlow: true, label: 'On Fire'),
  legendary(color: Color(0xFFFFD700), hasGlow: true, label: 'Legendary');

  final Color color;
  final bool hasGlow;
  final String label;

  const StreakTier({
    required this.color,
    required this.hasGlow,
    required this.label,
  });
}

/// Animated streak badge with tier label.
class StreakBadge extends StatelessWidget {
  final int streak;
  final bool compact;

  const StreakBadge({
    super.key,
    required this.streak,
    this.compact = false,
  });

  @override
  Widget build(BuildContext context) {
    if (streak == 0) {
      return const SizedBox.shrink();
    }

    final tier = _getStreakTier(streak);

    if (compact) {
      return Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        decoration: BoxDecoration(
          color: tier.color.withValues(alpha: 0.2),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: tier.color.withValues(alpha: 0.5)),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.local_fire_department, color: tier.color, size: 14),
            const SizedBox(width: 4),
            Text(
              '$streak',
              style: TextStyle(
                color: tier.color,
                fontSize: 12,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
      );
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            tier.color.withValues(alpha: 0.3),
            tier.color.withValues(alpha: 0.1),
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: tier.color.withValues(alpha: 0.5)),
        boxShadow: tier.hasGlow
            ? [
                BoxShadow(
                  color: tier.color.withValues(alpha: 0.2),
                  blurRadius: 8,
                  spreadRadius: 1,
                ),
              ]
            : null,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.local_fire_department, color: tier.color, size: 20),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                '$streak day${streak == 1 ? '' : 's'}',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                tier.label,
                style: TextStyle(
                  color: tier.color,
                  fontSize: 10,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  StreakTier _getStreakTier(int streak) {
    if (streak >= 100) {
      return StreakTier.legendary;
    } else if (streak >= 30) {
      return StreakTier.fire;
    } else if (streak >= 7) {
      return StreakTier.warm;
    } else {
      return StreakTier.spark;
    }
  }
}

/// Score indicator showing 0-100 score with trend arrow.
class ScoreIndicator extends StatelessWidget {
  final double score;
  final String? trend; // 'improving', 'stable', 'declining'
  final bool showLabel;

  const ScoreIndicator({
    super.key,
    required this.score,
    this.trend,
    this.showLabel = true,
  });

  @override
  Widget build(BuildContext context) {
    final color = _getScoreColor(score);
    final trendIcon = _getTrendIcon();

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 36,
          height: 36,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: color.withValues(alpha: 0.2),
            border: Border.all(color: color, width: 2),
          ),
          child: Center(
            child: Text(
              score.toInt().toString(),
              style: TextStyle(
                color: color,
                fontSize: 12,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
        if (trendIcon != null) ...[
          const SizedBox(width: 4),
          trendIcon,
        ],
        if (showLabel) ...[
          const SizedBox(width: 8),
          Text(
            'Score',
            style: TextStyle(
              color: Colors.grey[500],
              fontSize: 10,
            ),
          ),
        ],
      ],
    );
  }

  Color _getScoreColor(double score) {
    if (score >= 80) return Colors.green;
    if (score >= 60) return Colors.yellow.shade700;
    if (score >= 40) return Colors.orange;
    return Colors.red;
  }

  Widget? _getTrendIcon() {
    if (trend == null) return null;

    switch (trend) {
      case 'improving':
        return const Icon(Icons.trending_up, color: Colors.green, size: 16);
      case 'declining':
        return const Icon(Icons.trending_down, color: Colors.red, size: 16);
      case 'stable':
        return Icon(Icons.trending_flat, color: Colors.grey[500], size: 16);
      default:
        return null;
    }
  }
}
