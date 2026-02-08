import 'package:flutter/material.dart';

import '../models/mode.dart';

typedef ModeChanged = void Function(Mode mode);

class ModeToggle extends StatelessWidget {
  const ModeToggle({super.key, required this.mode, required this.onChanged});

  final Mode mode;
  final ModeChanged onChanged;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: theme.colorScheme.outline.withOpacity(0.2),
        ),
      ),
      child: Row(
        children: [
          _toggleItem(theme, Mode.bestHand, 'Best Hand'),
          _toggleItem(theme, Mode.headsUp, 'Heads-Up'),
          _toggleItem(theme, Mode.odds, 'Monte Carlo'),
        ],
      ),
    );
  }

  Widget _toggleItem(ThemeData theme, Mode target, String label) {
    final selected = mode == target;
    return Expanded(
      child: GestureDetector(
        onTap: () => onChanged(target),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.symmetric(vertical: 12),
          decoration: BoxDecoration(
            color: selected
                ? theme.colorScheme.primaryContainer
                : Colors.transparent,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Text(
            label,
            textAlign: TextAlign.center,
            style: theme.textTheme.labelLarge?.copyWith(
              color: selected
                  ? theme.colorScheme.onPrimaryContainer
                  : theme.colorScheme.onSurface.withOpacity(0.7),
            ),
          ),
        ),
      ),
    );
  }
}
