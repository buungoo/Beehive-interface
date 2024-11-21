import 'package:flutter/foundation.dart';
import 'package:flutter/cupertino.dart'; // Import cupertino package for CupertinoApp and CupertinoNavigationBar
import 'package:flutter/material.dart';
import '../widgets/shared.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';
import 'package:beehive/providers/beehive_user_provider.dart';

import 'dart:developer';
import 'dart:io';

class Camera extends StatefulWidget {
  const Camera({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _Camera();
}

class _Camera extends State<Camera> {
  Barcode? result;
  QRViewController? controller;
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');

  @override
  void reassemble() {
    super.reassemble();
    if (Platform.isAndroid) {
      controller!.pauseCamera();
    }
    controller!.resumeCamera();
  }

  @override
  void dispose() {
    controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SharedScaffold(
      context: context,
      appBar: getNavigationBar(
          context: context, title: "Camera", bgcolor: Color(0xFFf4991a)),
      body: Column(
        children: <Widget>[
          Expanded(
            flex: 5,
            child: _buildQrView(context),
          ),
          Expanded(
              child: FittedBox(
                  fit: BoxFit.contain,
                  child: Column(
                    children: <Widget>[
                      if (result != null)
                        Text(
                            'Barcode Type: ${describeEnum(result!.format)}   Data: ${result!.code}')
                      else
                        const Text('Scan a code'),
                    ],
                  )))
        ],
      ),
    );
  }

  Widget _buildQrView(BuildContext context) {
    // For this example we check how width or tall the device is and change the scanArea and overlay accordingly.
    var scanArea = (MediaQuery.of(context).size.width < 400 ||
            MediaQuery.of(context).size.height < 400)
        ? 150.0
        : 300.0;
    // To ensure the Scanner view is properly sizes after rotation
    // we need to listen for Flutter SizeChanged notification and update controller
    return QRView(
      key: qrKey,
      onQRViewCreated: _onQRViewCreated,
      overlay: QrScannerOverlayShape(
          borderColor: Colors.red,
          borderRadius: 10,
          borderLength: 30,
          borderWidth: 10,
          cutOutSize: scanArea),
      onPermissionSet: (ctrl, p) => _onPermissionSet(context, ctrl, p),
    );
  }

  void _onQRViewCreated(QRViewController controller) {
    setState(() {
      this.controller = controller;
    });
    controller.scannedDataStream.listen((scanData) async {
      print(scanData.format);
      if (scanData.format == BarcodeFormat.qrcode) {
        controller.pauseCamera();
        final token = scanData.code!;
        final response = await BeehiveUserProvider().addBeehive(token);

        //TODO: Let the user know
      }
      setState(() {
        result = scanData;
      });
    });
  }

  void _onPermissionSet(BuildContext context, QRViewController ctrl, bool p) {
    log('${DateTime.now().toIso8601String()}_onPermissionSet $p');
    if (!p) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('no Permission')),
      );
    }
  }
}

/**
 * body: Column(
    children: [
    Expanded(
    flex: 5,
    child: QRView(
    key: _qrKey,
    onQRViewCreated: _onQRViewCreated,
    ),
 */
