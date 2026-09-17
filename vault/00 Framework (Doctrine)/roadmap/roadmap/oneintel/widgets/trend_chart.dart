import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

/// Trend chart widget showing performance over time.
///
/// Uses fl_chart for smooth line graphs.
/// Add fl_chart to pubspec.yaml: fl_chart: ^0.68.0
class TrendChart extends StatelessWidget {
  final List<DataPoint> dataPoints;
  final int periodDays;
  final String? title;
  final Color lineColor;

  const TrendChart({
    super.key,
    required this.dataPoints,
    this.periodDays = 30,
    this.title,
    this.lineColor = Colors.green,
  });

  @override
  Widget build(BuildContext context) {
    if (dataPoints.isEmpty) {
      return _buildEmptyState();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (title != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Text(
              title!,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        SizedBox(
          height: 150,
          child: LineChart(
            LineChartData(
              gridData: FlGridData(
                show: true,
                drawVerticalLine: false,
                horizontalInterval: 25,
                getDrawingHorizontalLine: (value) => FlLine(
                  color: Colors.grey[800]!,
                  strokeWidth: 1,
                ),
              ),
              titlesData: FlTitlesData(
                leftTitles: AxisTitles(
                  sideTitles: SideTitles(
                    showTitles: true,
                    reservedSize: 30,
                    interval: 25,
                    getTitlesWidget: (value, meta) {
                      return Text(
                        '${value.toInt()}%',
                        style: TextStyle(
                          color: Colors.grey[600],
                          fontSize: 10,
                        ),
                      );
                    },
                  ),
                ),
                bottomTitles: AxisTitles(
                  sideTitles: SideTitles(
                    showTitles: true,
                    reservedSize: 20,
                    interval: _getBottomInterval(),
                    getTitlesWidget: _bottomTitleWidget,
                  ),
                ),
                topTitles: const AxisTitles(
                  sideTitles: SideTitles(showTitles: false),
                ),
                rightTitles: const AxisTitles(
                  sideTitles: SideTitles(showTitles: false),
                ),
              ),
              borderData: FlBorderData(show: false),
              minX: 0,
              maxX: (dataPoints.length - 1).toDouble(),
              minY: 0,
              maxY: 100,
              lineBarsData: [
                LineChartBarData(
                  spots: dataPoints.asMap().entries.map((e) {
                    return FlSpot(e.key.toDouble(), e.value.value);
                  }).toList(),
                  isCurved: true,
                  curveSmoothness: 0.3,
                  color: lineColor,
                  barWidth: 3,
                  isStrokeCapRound: true,
                  dotData: FlDotData(
                    show: dataPoints.length <= 14,
                    getDotPainter: (spot, percent, barData, index) {
                      return FlDotCirclePainter(
                        radius: 4,
                        color: lineColor,
                        strokeWidth: 2,
                        strokeColor: Colors.white,
                      );
                    },
                  ),
                  belowBarData: BarAreaData(
                    show: true,
                    gradient: LinearGradient(
                      colors: [
                        lineColor.withValues(alpha: 0.3),
                        lineColor.withValues(alpha: 0.0),
                      ],
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                    ),
                  ),
                ),
              ],
              lineTouchData: LineTouchData(
                enabled: true,
                touchTooltipData: LineTouchTooltipData(
                  tooltipBgColor: Colors.grey[850]!,
                  tooltipRoundedRadius: 8,
                  getTooltipItems: (touchedSpots) {
                    return touchedSpots.map((spot) {
                      final dataPoint = dataPoints[spot.x.toInt()];
                      return LineTooltipItem(
                        '${dataPoint.value.toInt()}%\n${_formatDate(dataPoint.date)}',
                        const TextStyle(
                          color: Colors.white,
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                        ),
                      );
                    }).toList();
                  },
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildEmptyState() {
    return Container(
      height: 150,
      alignment: Alignment.center,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.show_chart, color: Colors.grey[700], size: 32),
          const SizedBox(height: 8),
          Text(
            'No data yet',
            style: TextStyle(color: Colors.grey[600], fontSize: 12),
          ),
        ],
      ),
    );
  }

  double _getBottomInterval() {
    if (dataPoints.length <= 7) return 1;
    if (dataPoints.length <= 14) return 2;
    if (dataPoints.length <= 30) return 7;
    return 14;
  }

  Widget _bottomTitleWidget(double value, TitleMeta meta) {
    final index = value.toInt();
    if (index >= dataPoints.length) return const SizedBox.shrink();

    final date = dataPoints[index].date;
    return Padding(
      padding: const EdgeInsets.only(top: 4),
      child: Text(
        '${date.month}/${date.day}',
        style: TextStyle(
          color: Colors.grey[600],
          fontSize: 10,
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    const months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
    ];
    return '${months[date.month - 1]} ${date.day}';
  }
}

/// Data point for trend chart.
class DataPoint {
  final DateTime date;
  final double value; // 0-100

  const DataPoint({
    required this.date,
    required this.value,
  });
}

/// Period selector for trend chart timeframes.
class PeriodSelector extends StatelessWidget {
  final int selectedDays;
  final Function(int) onPeriodChanged;
  final List<int> periods;

  const PeriodSelector({
    super.key,
    required this.selectedDays,
    required this.onPeriodChanged,
    this.periods = const [7, 30, 90],
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: periods.map((days) {
        return Padding(
          padding: const EdgeInsets.only(right: 8),
          child: _buildPeriodChip(days, '${days}d'),
        );
      }).toList(),
    );
  }

  Widget _buildPeriodChip(int days, String label) {
    final isSelected = selectedDays == days;

    return GestureDetector(
      onTap: () => onPeriodChanged(days),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(
          color: isSelected ? Colors.red : Colors.grey[850],
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isSelected ? Colors.red : Colors.grey[700]!,
          ),
        ),
        child: Text(
          label,
          style: TextStyle(
            color: isSelected ? Colors.white : Colors.grey[400],
            fontSize: 12,
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}

/// Completion rate bar chart (vertical bars for each day/week).
class CompletionBars extends StatelessWidget {
  final Map<String, double> rates; // {'Mon': 85.0, 'Tue': 70.0, ...}

  const CompletionBars({
    super.key,
    required this.rates,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: rates.entries.map((entry) {
        return _buildBar(entry.key, entry.value);
      }).toList(),
    );
  }

  Widget _buildBar(String label, double value) {
    final color = _getBarColor(value);

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 24,
          height: 60,
          decoration: BoxDecoration(
            color: Colors.grey[850],
            borderRadius: BorderRadius.circular(4),
          ),
          child: Stack(
            alignment: Alignment.bottomCenter,
            children: [
              Container(
                width: 24,
                height: (value / 100) * 60,
                decoration: BoxDecoration(
                  color: color,
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: Colors.grey[600],
            fontSize: 10,
          ),
        ),
      ],
    );
  }

  Color _getBarColor(double value) {
    if (value >= 80) return Colors.green;
    if (value >= 60) return Colors.yellow.shade700;
    if (value >= 40) return Colors.orange;
    return Colors.red[400]!;
  }
}
