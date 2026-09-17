import 'package:flutter/material.dart';

/// Mini heatmap showing last N days of completion.
///
/// GitHub-style contribution squares in a compact row.
class MiniHeatmap extends StatelessWidget {
  final Set<DateTime> completionDates;
  final int days;
  final double squareSize;

  const MiniHeatmap({
    super.key,
    required this.completionDates,
    this.days = 7,
    this.squareSize = 16,
  });

  @override
  Widget build(BuildContext context) {
    final today = DateTime.now();
    final daysList = List.generate(
      days,
      (i) => today.subtract(Duration(days: days - 1 - i)),
    );

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: daysList.map((day) {
        final isCompleted = _isDateCompleted(day);
        final isToday = _isSameDay(day, today);

        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 2),
          child: Container(
            width: squareSize,
            height: squareSize,
            decoration: BoxDecoration(
              color: isCompleted ? Colors.green[600] : Colors.grey[800],
              borderRadius: BorderRadius.circular(3),
              border: isToday
                  ? Border.all(color: Colors.red, width: 1)
                  : null,
            ),
          ),
        );
      }).toList(),
    );
  }

  bool _isDateCompleted(DateTime day) {
    return completionDates.any((d) => _isSameDay(d, day));
  }

  bool _isSameDay(DateTime a, DateTime b) {
    return a.year == b.year && a.month == b.month && a.day == b.day;
  }
}

/// Weekly heatmap showing M-T-W-T-F-S-S grid.
///
/// Shows completion status for current week with frequency awareness.
class WeeklyHeatmap extends StatelessWidget {
  final Set<DateTime> completionDates;
  final List<int> frequencyDays; // ISO weekdays (1=Mon, 7=Sun)

  const WeeklyHeatmap({
    super.key,
    required this.completionDates,
    this.frequencyDays = const [1, 2, 3, 4, 5, 6, 7],
  });

  @override
  Widget build(BuildContext context) {
    final today = DateTime.now();
    final startOfWeek = today.subtract(Duration(days: today.weekday - 1));

    const dayLabels = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: List.generate(7, (i) {
        final day = startOfWeek.add(Duration(days: i));
        final isRequired = frequencyDays.contains(i + 1);
        final isCompleted = _isDateCompleted(day);
        final isToday = i + 1 == today.weekday;
        final isFuture = day.isAfter(today);

        Color backgroundColor;
        Color textColor;

        if (isFuture) {
          backgroundColor = Colors.grey[900]!;
          textColor = Colors.grey[700]!;
        } else if (isCompleted) {
          backgroundColor = Colors.green[600]!;
          textColor = Colors.white;
        } else if (isRequired && !isFuture) {
          backgroundColor = Colors.red[900]!.withValues(alpha: 0.5);
          textColor = Colors.red[300]!;
        } else {
          backgroundColor = Colors.grey[850]!;
          textColor = Colors.grey[500]!;
        }

        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              dayLabels[i],
              style: TextStyle(
                color: isRequired ? Colors.white70 : Colors.grey[600],
                fontSize: 10,
                fontWeight: isRequired ? FontWeight.w600 : FontWeight.normal,
              ),
            ),
            const SizedBox(height: 4),
            Container(
              width: 32,
              height: 32,
              decoration: BoxDecoration(
                color: backgroundColor,
                borderRadius: BorderRadius.circular(6),
                border: isToday
                    ? Border.all(color: Colors.red, width: 2)
                    : null,
              ),
              child: Center(
                child: isCompleted
                    ? const Icon(Icons.check, size: 18, color: Colors.white)
                    : Text(
                        day.day.toString(),
                        style: TextStyle(color: textColor, fontSize: 12),
                      ),
              ),
            ),
          ],
        );
      }),
    );
  }

  bool _isDateCompleted(DateTime day) {
    return completionDates.any((d) =>
        d.year == day.year &&
        d.month == day.month &&
        d.day == day.day);
  }
}

/// Month heatmap grid (7 columns x 5-6 rows).
class MonthHeatmap extends StatelessWidget {
  final Set<DateTime> completionDates;
  final DateTime? focusedMonth;
  final Function(DateTime)? onDayTap;

  const MonthHeatmap({
    super.key,
    required this.completionDates,
    this.focusedMonth,
    this.onDayTap,
  });

  @override
  Widget build(BuildContext context) {
    final month = focusedMonth ?? DateTime.now();
    final firstDay = DateTime(month.year, month.month, 1);
    final lastDay = DateTime(month.year, month.month + 1, 0);
    final today = DateTime.now();

    // Get starting weekday (0 = Monday in our grid)
    final startPadding = firstDay.weekday - 1;
    final totalDays = lastDay.day;
    final totalCells = startPadding + totalDays;
    final rows = (totalCells / 7).ceil();

    return Column(
      children: [
        // Day labels
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: ['M', 'T', 'W', 'T', 'F', 'S', 'S']
              .map((d) => SizedBox(
                    width: 32,
                    child: Center(
                      child: Text(
                        d,
                        style: TextStyle(
                          color: Colors.grey[600],
                          fontSize: 10,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ))
              .toList(),
        ),
        const SizedBox(height: 8),
        // Calendar grid
        ...List.generate(rows, (row) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 4),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: List.generate(7, (col) {
                final cellIndex = row * 7 + col;
                final dayNum = cellIndex - startPadding + 1;

                if (dayNum < 1 || dayNum > totalDays) {
                  return const SizedBox(width: 32, height: 32);
                }

                final day = DateTime(month.year, month.month, dayNum);
                final isCompleted = _isDateCompleted(day);
                final isToday = _isSameDay(day, today);
                final isFuture = day.isAfter(today);

                return GestureDetector(
                  onTap: onDayTap != null ? () => onDayTap!(day) : null,
                  child: Container(
                    width: 32,
                    height: 32,
                    decoration: BoxDecoration(
                      color: isFuture
                          ? Colors.grey[900]
                          : isCompleted
                              ? _getCompletionColor(day, today)
                              : Colors.grey[850],
                      borderRadius: BorderRadius.circular(6),
                      border: isToday
                          ? Border.all(color: Colors.red, width: 2)
                          : null,
                    ),
                    child: Center(
                      child: Text(
                        dayNum.toString(),
                        style: TextStyle(
                          color: isFuture
                              ? Colors.grey[700]
                              : isCompleted
                                  ? Colors.white
                                  : Colors.white70,
                          fontSize: 12,
                          fontWeight:
                              isCompleted ? FontWeight.w600 : FontWeight.normal,
                        ),
                      ),
                    ),
                  ),
                );
              }),
            ),
          );
        }),
      ],
    );
  }

  bool _isDateCompleted(DateTime day) {
    return completionDates.any((d) => _isSameDay(d, day));
  }

  bool _isSameDay(DateTime a, DateTime b) {
    return a.year == b.year && a.month == b.month && a.day == b.day;
  }

  Color _getCompletionColor(DateTime day, DateTime today) {
    final daysAgo = today.difference(day).inDays;
    if (daysAgo < 7) return Colors.green[700]!;
    if (daysAgo < 30) return Colors.green[500]!;
    return Colors.green[300]!;
  }
}
