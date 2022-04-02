import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:front/constants.dart';
import 'package:front/controller/main_controller.dart';
import 'package:get/get.dart';
import 'package:http/http.dart' as http;

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return GetMaterialApp(
      title: 'Smart Drug Box',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(title: 'Smart Drug Box'),
      initialBinding: BindingsBuilder(() {
        Get.put(MainController());
      }),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);
  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final hText = TextEditingController(text: "00");
  final mText = TextEditingController(text: "00");

  void _incrementCounter(MainController ctrl) async {
    var apiUrl = Uri.base.path[Uri.base.path.length - 1] == '/'
        ? '${Uri.base.path}api/v1/alarm'
        : '${Uri.base.path}/api/v1/alarm';

    var body =
        jsonEncode({"h": int.parse(hText.text), "m": int.parse(mText.text)});
    await http.post(
        Uri.http(
          serverAddr,
          apiUrl,
        ),
        body: body);

    ctrl.loadData();
  }

  @override
  void dispose() {
    hText.dispose();
    mText.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            const Text('Current Alarms', style: TextStyle(fontSize: 26)),
            GetBuilder<MainController>(builder: (ctrl) {
              return Text(ctrl.alarmList.value,
                  style: const TextStyle(fontSize: 22));
            }),
            const Text('알람 추가', style: TextStyle(fontSize: 26)),
            Row(
              mainAxisSize: MainAxisSize.min,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(
                  width: 80,
                  child: Center(
                    child: TextField(
                      inputFormatters: [LengthLimitingTextInputFormatter(2)],
                      decoration:
                          const InputDecoration(border: OutlineInputBorder()),
                      textAlign: TextAlign.center,
                      controller: hText,
                      style: Theme.of(context).textTheme.headline4,
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                Text(
                  ':',
                  style: Theme.of(context).textTheme.headline4,
                ),
                const SizedBox(width: 10),
                SizedBox(
                  width: 80,
                  child: Center(
                    child: TextField(
                      inputFormatters: [LengthLimitingTextInputFormatter(2)],
                      decoration:
                          const InputDecoration(border: OutlineInputBorder()),
                      textAlign: TextAlign.center,
                      controller: mText,
                      style: Theme.of(context).textTheme.headline4,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
      floatingActionButton: GetBuilder<MainController>(builder: (ctrl) {
        return FloatingActionButton(
          onPressed: () {
            _incrementCounter(ctrl);
          },
          tooltip: 'Add',
          child: const Icon(Icons.add),
        );
      }), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
