import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() => runApp(const MyApp());

class Note {
  final int id;
  final String title;

  Note({required this.id, required this.title});

  factory Note.fromJson(Map<String, dynamic> json) {
    return Note(id: json['id'], title: json['title']);
  }
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  List<Note> notes = [];
  final TextEditingController controller = TextEditingController();

  final String baseUrl = "http://localhost:8000";

  Future<void> fetchNotes() async {
    final res = await http.get(Uri.parse("$baseUrl/notes"));
    if (res.statusCode == 200) {
      final List<dynamic> data = jsonDecode(res.body);
      setState(() {
        notes = data.map((e) => Note.fromJson(e)).toList();
      });
    }
  }

  Future<void> addNote() async {
    final res = await http.post(
      Uri.parse("$baseUrl/notes"),
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"title": controller.text}),
    );
    if (res.statusCode == 200) {
      controller.clear();
      fetchNotes();
    }
  }

  @override
  void initState() {
    super.initState();
    fetchNotes();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(title: const Text("Notes App")),
        body: Column(
          children: [
            TextField(
              controller: controller,
              decoration: InputDecoration(
                hintText: "Enter note",
                suffixIcon: IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: addNote,
                ),
              ),
            ),
            Expanded(
              child: RefreshIndicator(
                onRefresh: fetchNotes,
                child: ListView.builder(
                  itemCount: notes.length,
                  itemBuilder: (context, i) =>
                      ListTile(title: Text(notes[i].title)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
