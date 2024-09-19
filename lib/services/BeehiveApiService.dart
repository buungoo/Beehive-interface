import '../models/beehive.dart';

class BeehiveApi {
  BeehiveApi() {}

  /**
   * Return an array of all beehives that are available for the user
   */
  Future<List<Beehive>> GetHives() {
    return Future.delayed(
        const Duration(seconds: 10),
        () => [
              Beehive(id: "1", name: "Beehive 1"),
              Beehive(id: "2", name: "Beehive 2"),
            ]);
  }

/**
 * Return a single beehive by its id
 */
/*GetHive(int id) {
    return}*/
}
