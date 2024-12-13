// Import cupertino package for CupertinoApp and CupertinoNavigationBar
import 'package:flutter/material.dart';
import '../widgets/shared.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';
import 'package:beehive/providers/beehive_user_provider.dart';
// GoRouter for navigation

import 'dart:developer';
import 'dart:io';

class Camera extends StatefulWidget {
  const Camera({super.key});

  @override
  State<StatefulWidget> createState() => _Camera();
}

class _Camera extends State<Camera> {
  Barcode? result;
  QRViewController? controller;
  bool isLoading = false;
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
          context: context, title: "Camera", bgcolor: const Color(0xFFf4991a)),
      body: Column(
        children: <Widget>[
          Expanded(
            flex: 5,
            child: _buildQrView(context),
          ),
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
      if (isLoading) return;

      setState(() {
        isLoading = true;
      });

      if (scanData.format == BarcodeFormat.qrcode) {
        controller.pauseCamera();
        final token = scanData.code!;
        final response = await BeehiveUserProvider().addBeehive(token!);

        if (response) {
          showDialog(
            context: context,
            builder: (BuildContext context) {
              return AlertDialog(
                title: const Text('QR Code Scanned'),
                content: const Text('Beehive added successfully!'),
                actions: <Widget>[
                  TextButton(
                    child: const Text('OK'),
                    onPressed: () {
                      Navigator.pop(context);
                      Navigator.pop(context, true);
                    },
                  ),
                ],
              );
            },
          );
        } else {
          showDialog(
            context: context,
            builder: (BuildContext context) {
              return AlertDialog(
                title: const Text('QR Code Scanned'),
                content: const Text('Failed to add beehive. Please try again.'),
                actions: <Widget>[
                  TextButton(
                    child: const Text('OK'),
                    onPressed: () async {
                      controller.resumeCamera();
                      Navigator.pop(context);
                    },
                  ),
                ],
              );
            },
          );
        }

        // Show an alert dialog to let the user know the result
      }
      setState(() {
        result = scanData;
        isLoading = false;
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
