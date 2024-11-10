class Beehive {
  final String id;
  final String name;

  Beehive({required this.id, required this.name});

  factory Beehive.fromJson(Map<String, dynamic> json) {
    return Beehive(
      id: json['Id'].toString(),
      name: json['Name'],
    );
  }
}
