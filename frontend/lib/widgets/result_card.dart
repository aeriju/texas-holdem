import 'package:flutter/material.dart';

import '../models/mode.dart';

class ResultCard extends StatelessWidget {
  const ResultCard({super.key, required this.mode, required this.result});

  final Mode mode;
  final Map<String, dynamic>? result;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: theme.colorScheme.outline.withOpacity(0.2),
        ),
      ),
      child: result == null
          ? Text(
              'Run a request to see output here.',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withOpacity(0.6),
              ),
            )
          : Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                ..._buildFormattedResult(theme),
                const SizedBox(height: 12),
              ],
            ),
    );
  }

  List<Widget> _buildFormattedResult(ThemeData theme) {
    if (result == null) return [];
    switch (mode) {
      case Mode.bestHand:
        return _buildBestHandResult(theme, result!);
      case Mode.headsUp:
        return _buildHeadsUpResult(theme, result!);
      case Mode.odds:
        return _buildOddsResult(theme, result!);
    }
  }

  List<Widget> _buildBestHandResult(ThemeData theme, Map<String, dynamic> data) {
    final category = data['category']?.toString() ?? 'Unknown';
    final bestHand = (data['bestHand'] as List<dynamic>? ?? [])
        .map((e) => e.toString())
        .toList();
    return [
      Text(
        category.toUpperCase(),
        style: theme.textTheme.titleMedium?.copyWith(
          color: theme.colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
      const SizedBox(height: 6),
      Wrap(
        spacing: 8,
        runSpacing: 8,
        children: bestHand.map(_cardChip).toList(),
      ),
    ];
  }

  List<Widget> _buildHeadsUpResult(ThemeData theme, Map<String, dynamic> data) {
    final winner = data['winner']?.toString() ?? 'tie';
    final outcome = data['outcome']?.toString() ?? 'tie';
    final hand1 = data['hand1'] as Map<String, dynamic>? ?? {};
    final hand2 = data['hand2'] as Map<String, dynamic>? ?? {};
    final h1Category = hand1['category']?.toString() ?? 'Unknown';
    final h2Category = hand2['category']?.toString() ?? 'Unknown';
    final h1Best = (hand1['bestHand'] as List<dynamic>? ?? [])
        .map((e) => e.toString())
        .toList();
    final h2Best = (hand2['bestHand'] as List<dynamic>? ?? [])
        .map((e) => e.toString())
        .toList();

    return [
      Text(
        outcome.toUpperCase(),
        style: theme.textTheme.titleMedium?.copyWith(
          color: theme.colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
      const SizedBox(height: 8),
      Text('Player 1: $h1Category',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: winner == 'hand1'
                ? theme.colorScheme.primary
                : theme.colorScheme.onSurface.withOpacity(0.8),
          )),
      const SizedBox(height: 6),
      Wrap(spacing: 8, runSpacing: 8, children: h1Best.map(_cardChip).toList()),
      const SizedBox(height: 12),
      Text('Player 2: $h2Category',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: winner == 'hand2'
                ? theme.colorScheme.primary
                : theme.colorScheme.onSurface.withOpacity(0.8),
          )),
      const SizedBox(height: 6),
      Wrap(spacing: 8, runSpacing: 8, children: h2Best.map(_cardChip).toList()),
    ];
  }

  List<Widget> _buildOddsResult(ThemeData theme, Map<String, dynamic> data) {
    final prob = (data['winProbability'] as num?)?.toDouble();
    final pct = prob == null ? 'N/A' : '${(prob * 100).toStringAsFixed(2)}%';
    return [
      Text(
        'Win Probability',
        style: theme.textTheme.titleMedium?.copyWith(
          color: theme.colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
      const SizedBox(height: 6),
      Text(
        pct,
        style: theme.textTheme.headlineMedium?.copyWith(
          color: theme.colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
    ];
  }

  Widget _cardChip(String card) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: const Color(0xFF0F172A),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.white24),
      ),
      child: Text(
        card,
        style: const TextStyle(
          color: Colors.white,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.5,
        ),
      ),
    );
  }
}
