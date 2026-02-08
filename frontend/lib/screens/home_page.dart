import 'package:flutter/material.dart';

import '../models/mode.dart';
import '../services/api_client.dart';
import '../widgets/inputs.dart';
import '../widgets/mode_toggle.dart';
import '../widgets/result_card.dart';

class HoldemPage extends StatefulWidget {
  const HoldemPage({super.key});

  @override
  State<HoldemPage> createState() => _HoldemPageState();
}

class _HoldemPageState extends State<HoldemPage> {
  Mode _mode = Mode.bestHand;
  bool _loading = false;
  String? _error;
  Map<String, dynamic>? _result;

  final ApiClient _api = ApiClient();

  final _bestHole = TextEditingController(text: 'HA HK');
  final _bestCommunity = TextEditingController(text: 'C2 D3 S4 H5 C6');

  final _hu1Hole = TextEditingController(text: 'HA HK');
  final _hu1Community = TextEditingController(text: 'C2 D3 S4 H5 C6');
  final _hu2Hole = TextEditingController(text: 'S9 S8');
  final _hu2Community = TextEditingController(text: 'C2 D3 S4 H5 C6');

  final _oddsHole = TextEditingController(text: 'HA HK');
  final _oddsCommunity = TextEditingController(text: 'C2 D3 S4');
  final _oddsPlayers = TextEditingController(text: '6');
  final _oddsSims = TextEditingController(text: '20000');

  @override
  void dispose() {
    _bestHole.dispose();
    _bestCommunity.dispose();
    _hu1Hole.dispose();
    _hu1Community.dispose();
    _hu2Hole.dispose();
    _hu2Community.dispose();
    _oddsHole.dispose();
    _oddsCommunity.dispose();
    _oddsPlayers.dispose();
    _oddsSims.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      body: Container(
        width: double.infinity,
        height: double.infinity,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              theme.colorScheme.surface,
              theme.colorScheme.surface.withOpacity(0.95),
              const Color(0xFF0B1220),
            ],
          ),
        ),
        child: SafeArea(
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24),
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 520),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Text(
                      'Texas Hold\'em',
                      style: theme.textTheme.headlineLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: theme.colorScheme.primary,
                        letterSpacing: -0.5,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Evaluate hands, compare matchups, or estimate win probability',
                      style: theme.textTheme.bodyLarge?.copyWith(
                        color: theme.colorScheme.onSurface.withOpacity(0.7),
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 32),
                    Container(
                      padding: const EdgeInsets.all(24),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainerHighest
                            .withOpacity(0.5),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: theme.colorScheme.outline.withOpacity(0.2),
                        ),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          ModeToggle(
                            mode: _mode,
                            onChanged: (mode) => setState(() {
                              _mode = mode;
                              _result = null;
                              _error = null;
                            }),
                          ),
                          const SizedBox(height: 20),
                          ModeInputs(
                            mode: _mode,
                            bestHole: _bestHole,
                            bestCommunity: _bestCommunity,
                            hu1Hole: _hu1Hole,
                            hu1Community: _hu1Community,
                            hu2Hole: _hu2Hole,
                            hu2Community: _hu2Community,
                            oddsHole: _oddsHole,
                            oddsCommunity: _oddsCommunity,
                            oddsPlayers: _oddsPlayers,
                            oddsSims: _oddsSims,
                            fieldBuilder: _field,
                          ),
                          const SizedBox(height: 20),
                          FilledButton(
                            onPressed: _loading ? null : _submit,
                            style: FilledButton.styleFrom(
                              padding: const EdgeInsets.symmetric(vertical: 18),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: _loading
                                ? SizedBox(
                                    height: 24,
                                    width: 24,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      color: theme.colorScheme.onPrimary,
                                    ),
                                  )
                                : const Text('Run'),
                          ),
                          if (_error != null) ...[
                            const SizedBox(height: 16),
                            Text(
                              _error!,
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: theme.colorScheme.error,
                              ),
                            ),
                          ],
                          const SizedBox(height: 20),
                          ResultCard(mode: _mode, result: _result),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _field(TextEditingController controller,
      {required String label, required String hint, TextInputType? keyboardType}) {
    final theme = Theme.of(context);
    return TextField(
      controller: controller,
      keyboardType: keyboardType,
      style: theme.textTheme.titleMedium,
      decoration: InputDecoration(
        labelText: label,
        hintText: hint,
        filled: true,
        fillColor: theme.colorScheme.surface,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 20,
          vertical: 16,
        ),
      ),
    );
  }

  Future<void> _submit() async {
    setState(() {
      _loading = true;
      _error = null;
      _result = null;
    });

    try {
      final payload = _buildPayload();
      final data = await _api.post(_mode, payload);
      setState(() {
        _result = data;
      });
    } catch (e) {
      setState(() {
        _error = 'Request failed: $e';
        _result = null;
      });
    } finally {
      setState(() {
        _loading = false;
      });
    }
  }

  Map<String, dynamic> _buildPayload() {
    switch (_mode) {
      case Mode.bestHand:
        return {
          'hole': _splitCards(_bestHole.text),
          'community': _splitCards(_bestCommunity.text),
        };
      case Mode.headsUp:
        return {
          'hand1': {
            'hole': _splitCards(_hu1Hole.text),
            'community': _splitCards(_hu1Community.text),
          },
          'hand2': {
            'hole': _splitCards(_hu2Hole.text),
            'community': _splitCards(_hu2Community.text),
          },
        };
      case Mode.odds:
        return {
          'hole': _splitCards(_oddsHole.text),
          'community': _splitCards(_oddsCommunity.text),
          'players': int.tryParse(_oddsPlayers.text.trim()) ?? 0,
          'simulations': int.tryParse(_oddsSims.text.trim()) ?? 0,
        };
    }
  }

  List<String> _splitCards(String text) {
    final cleaned = text.trim().replaceAll(',', ' ');
    if (cleaned.isEmpty) return [];
    return cleaned.split(RegExp(r'\s+')).where((e) => e.isNotEmpty).toList();
  }
}
