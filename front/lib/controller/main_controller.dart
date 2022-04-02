import 'dart:convert';

import 'package:front/constants.dart';
import 'package:get/get.dart';
import 'package:http/http.dart' as http;

class MainController extends GetxController {
  var alarmList = "".obs;

  MainController() {
    loadData();
  }

  void loadData() async {
    var apiUrl = Uri.base.path[Uri.base.path.length - 1] == '/'
        ? '${Uri.base.path}api/v1/alarm'
        : '${Uri.base.path}/api/v1/alarm';
    var response = await http.get(Uri.http(serverAddr, apiUrl));
    alarmList.value = response.body;

    update();
  }
}
