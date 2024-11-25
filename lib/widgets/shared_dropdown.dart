import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:beehive/utils/helpers.dart';

class SharedDropdownMenu extends StatefulWidget {
  final List<String> itemList;
  final ValueChanged<String> onItemChanged;

  const SharedDropdownMenu(
      {super.key, required this.itemList, required this.onItemChanged});

  @override
  State<SharedDropdownMenu> createState() => _SharedDropdownMenuState();
}

class _SharedDropdownMenuState extends State<SharedDropdownMenu> {
  int _selectedFruit = 0;

  @override
  void dispose() {
    // Perform any necessary cleanup here
    super.dispose();
  }

  void _showDialog(Widget child) {
    showCupertinoModalPopup<void>(
      context: context,
      builder: (BuildContext context) => Container(
        height: 216,
        padding: const EdgeInsets.only(top: 6.0),
        // The Bottom margin is provided to align the popup above the system navigation bar.
        margin: EdgeInsets.only(
          bottom: MediaQuery.of(context).viewInsets.bottom,
        ),
        // Provide a background color for the popup.
        color: CupertinoColors.systemBackground.resolveFrom(context),
        // Use a SafeArea widget to avoid system overlaps.
        child: SafeArea(
          top: false,
          child: child,
        ),
      ),
    );
  }

  Widget _iosBuild(BuildContext context) {
    const double _kItemExtent = 32.0;
    List<String> list = widget.itemList;
    String dropdownValue = list.first;
    return CupertinoButton(
      padding: EdgeInsets.zero,
      // Display a CupertinoPicker with list of fruits.
      onPressed: () => _showDialog(
        CupertinoPicker(
          magnification: 1.22,
          squeeze: 1.2,
          useMagnifier: true,
          itemExtent: _kItemExtent,
          // This sets the initial item.
          scrollController: FixedExtentScrollController(
            initialItem: _selectedFruit,
          ),
          // This is called when selected item is changed.
          onSelectedItemChanged: (int selectedItem) {
            if (mounted) {
              setState(() {
                _selectedFruit = selectedItem;
                widget.onItemChanged(list[selectedItem]);
              });
            }
          },
          children: List<Widget>.generate(list.length, (int index) {
            return Center(child: Text(list[index]));
          }),
        ),
      ),
      // This displays the selected fruit name.
      child: Text(
        list[_selectedFruit],
        style: const TextStyle(
          fontSize: 22.0,
        ),
      ),
    );
  }

  @override
  Widget _androidBuild(BuildContext context) {
    List<String> list = widget.itemList;
    String dropdownValue = list.first;
    return DropdownMenu<String>(
      initialSelection: list.first,
      onSelected: (String? value) {
        // This is called when the user selects an item.
        setState(() {
          dropdownValue = value!;
        });
      },
      dropdownMenuEntries: list.map<DropdownMenuEntry<String>>((String value) {
        return DropdownMenuEntry<String>(value: value, label: value);
      }).toList(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return isIOS(context) ? _iosBuild(context) : _androidBuild(context);
  }
}
