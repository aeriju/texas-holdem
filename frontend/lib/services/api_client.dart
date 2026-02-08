import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/mode.dart';

class ApiClient {
  ApiClient({String? baseUrl})
      : _baseUrl = baseUrl ?? const String.fromEnvironment(
          'REST_BACKEND_URL',
          defaultValue: 'http://localhost:8080',
        );

  final String _baseUrl;

  Future<Map<String, dynamic>> post(Mode mode, Map<String, dynamic> payload) async {
    final uri = _buildUri(mode);
    final res = await http.post(
      uri,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(payload),
    );

    if (res.statusCode >= 400) {
      throw Exception('Server error (${res.statusCode}): ${res.body}');
    }

    return jsonDecode(res.body) as Map<String, dynamic>;
  }

  Uri _buildUri(Mode mode) {
    final cleaned = _baseUrl.trim().replaceAll(RegExp(r'/+$'), '');
    switch (mode) {
      case Mode.bestHand:
        return Uri.parse('$cleaned/api/v1/best-hand');
      case Mode.headsUp:
        return Uri.parse('$cleaned/api/v1/heads-up');
      case Mode.odds:
        return Uri.parse('$cleaned/api/v1/odds');
    }
  }
}
