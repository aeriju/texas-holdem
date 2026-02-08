import 'package:flutter/material.dart';

import 'screens/home_page.dart';

class HoldemApp extends StatelessWidget {
  const HoldemApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        useMaterial3: true,
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF0EA5E9),
          brightness: Brightness.dark,
          primary: const Color(0xFF38BDF8),
          secondary: const Color(0xFF7DD3FC),
        ),
      ),
      home: const HoldemPage(),
    );
  }
}
