import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Route | ui/home', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let route = this.owner.lookup('route:ui/home');
    assert.ok(route);
  });
});
