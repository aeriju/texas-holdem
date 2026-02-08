import 'package:flutter/material.dart';

import '../models/mode.dart';

typedef PayloadBuilder = Map<String, dynamic> Function();

typedef TextFieldBuilder = Widget Function(
  TextEditingController controller, {
  required String label,
  required String hint,
  TextInputType? keyboardType,
});

class ModeInputs extends StatelessWidget {
  const ModeInputs({
    super.key,
    required this.mode,
    required this.bestHole,
    required this.bestCommunity,
    required this.hu1Hole,
    required this.hu1Community,
    required this.hu2Hole,
    required this.hu2Community,
    required this.oddsHole,
    required this.oddsCommunity,
    required this.oddsPlayers,
    required this.oddsSims,
    required this.fieldBuilder,
  });

  final Mode mode;
  final TextEditingController bestHole;
  final TextEditingController bestCommunity;
  final TextEditingController hu1Hole;
  final TextEditingController hu1Community;
  final TextEditingController hu2Hole;
  final TextEditingController hu2Community;
  final TextEditingController oddsHole;
  final TextEditingController oddsCommunity;
  final TextEditingController oddsPlayers;
  final TextEditingController oddsSims;
  final TextFieldBuilder fieldBuilder;

  @override
  Widget build(BuildContext context) {
    switch (mode) {
      case Mode.bestHand:
        return Column(
          children: [
            fieldBuilder(bestHole, label: 'Hole Cards', hint: 'HA HK'),
            const SizedBox(height: 16),
            fieldBuilder(bestCommunity,
                label: 'Community Cards', hint: 'C2 D3 S4 H5 C6'),
          ],
        );
      case Mode.headsUp:
        final theme = Theme.of(context);
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Player 1', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            fieldBuilder(hu1Hole, label: 'Hole Cards', hint: 'HA HK'),
            const SizedBox(height: 12),
            fieldBuilder(hu1Community,
                label: 'Community Cards', hint: 'C2 D3 S4 H5 C6'),
            const SizedBox(height: 16),
            Text('Player 2', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            fieldBuilder(hu2Hole, label: 'Hole Cards', hint: 'S9 S8'),
            const SizedBox(height: 12),
            fieldBuilder(hu2Community,
                label: 'Community Cards', hint: 'C2 D3 S4 H5 C6'),
          ],
        );
      case Mode.odds:
        return Column(
          children: [
            fieldBuilder(oddsHole, label: 'Hole Cards', hint: 'HA HK'),
            const SizedBox(height: 12),
            fieldBuilder(oddsCommunity,
                label: 'Community Cards', hint: 'C2 D3 S4'),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: fieldBuilder(
                    oddsPlayers,
                    label: 'Players',
                    hint: '6',
                    keyboardType: TextInputType.number,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: fieldBuilder(
                    oddsSims,
                    label: 'Simulations',
                    hint: '20000',
                    keyboardType: TextInputType.number,
                  ),
                ),
              ],
            ),
          ],
        );
    }
  }
}
